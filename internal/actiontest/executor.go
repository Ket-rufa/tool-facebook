package actiontest

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	graphQLEndpoint   = "https://www.facebook.com/api/graphql/"
	cometFriendlyName = "CometUFIFeedbackReactMutation"
	legacyFriendly    = "FeedbackReactMutation"
)

// Comet mutation currently used by web Facebook for reacting.
const cometReactionMutation = `mutation CometUFIFeedbackReactMutation($input: CometUFIFeedbackReactMutationInput!) {
  feedback_react(input: $input) {
    feedback {
      id
    }
  }
}`

// Fallback mutation for older schemas.
const legacyReactionMutation = `mutation FeedbackReactMutation($input: FeedbackReactInput!) {
  feedback_react(input: $input) {
    feedback {
      id
    }
  }
}`

type reactionStrategy struct {
	name         string
	friendlyName string
	query        string
	docID        string
}

func generateJazoest(fbDtsg string) string {
	sum := 0
	for _, c := range fbDtsg {
		sum += int(c)
	}
	return "2" + strconv.Itoa(sum)
}

func ExecuteFacebookReaction(cookie, fbDtsg, actorID, feedbackID, reactionID, configuredDocID string) (string, ActionErrorCode, error) {
	cookie = strings.TrimSpace(cookie)
	fbDtsg = strings.TrimSpace(fbDtsg)
	actorID = strings.TrimSpace(actorID)
	feedbackID = strings.TrimSpace(feedbackID)
	reactionID = strings.TrimSpace(reactionID)

	if cookie == "" || fbDtsg == "" {
		return "", ErrSessionMissing, errors.New("thieu cookie hoac fb_dtsg")
	}
	if actorID == "" {
		return "", ErrActorMismatch, errors.New("thieu actor_id (c_user)")
	}
	if reactionID == "" {
		return "", ErrMalformed, errors.New("thieu reaction_id")
	}
	if err := validateFeedbackID(feedbackID); err != nil {
		return "", ErrIDMismatch, err
	}

	variablesJSON, err := buildReactionVariables(actorID, feedbackID, reactionID)
	if err != nil {
		return "", ErrMalformed, fmt.Errorf("khong tao duoc GraphQL variables: %w", err)
	}

	strategies := buildReactionStrategies(configuredDocID)
	var lastErr error
	lastCode := ErrGraphqlFB

	for i, strategy := range strategies {
		log.Printf("[Executor] Attempt %d/%d strategy=%s friendly=%s use_doc_id=%t feedbackID=%q actorID=%q",
			i+1, len(strategies), strategy.name, strategy.friendlyName, strategy.docID != "", feedbackID, actorID)

		body, errCode, reqErr := executeReactionRequest(cookie, fbDtsg, actorID, variablesJSON, strategy)
		if reqErr == nil {
			return "Tha cam xuc thanh cong", "", nil
		}

		lastErr = reqErr
		if errCode != "" {
			lastCode = errCode
		}

		// If query shape is rejected, try next strategy before failing.
		if isIncorrectQuery(body) && i < len(strategies)-1 {
			log.Printf("[Executor] Incorrect Query with strategy=%s, fallback to next strategy", strategy.name)
			continue
		}
		return "", errCode, reqErr
	}

	if lastErr == nil {
		lastErr = errors.New("khong co strategy GraphQL nao thuc thi duoc")
	}
	return "", lastCode, lastErr
}

func validateFeedbackID(feedbackID string) error {
	if feedbackID == "" {
		return errors.New("feedback_id khong duoc de trong")
	}

	// Fast path for canonical prefix of base64("feedback:...").
	if strings.HasPrefix(feedbackID, "ZmVlZGJhY2s6") {
		return nil
	}

	decoded, err := decodeBase64(feedbackID)
	if err != nil {
		return errors.New("feedback_id khong dung base64 hop le")
	}
	if !strings.HasPrefix(decoded, "feedback:") {
		return errors.New("feedback_id phai la base64 cua chuoi bat dau bang 'feedback:'")
	}
	return nil
}

func decodeBase64(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("empty")
	}

	if b, err := base64.StdEncoding.DecodeString(raw); err == nil {
		return string(b), nil
	}

	padded := raw
	if m := len(raw) % 4; m != 0 {
		padded = raw + strings.Repeat("=", 4-m)
	}
	if b, err := base64.StdEncoding.DecodeString(padded); err == nil {
		return string(b), nil
	}

	if b, err := base64.RawStdEncoding.DecodeString(raw); err == nil {
		return string(b), nil
	}

	return "", errors.New("invalid base64")
}

func buildReactionVariables(actorID, feedbackID, reactionID string) (string, error) {
	nowMillis := time.Now().UnixMilli()
	attribution := fmt.Sprintf(
		"CometSinglePostDialogRoot.react,comet.post.single_dialog,via_cold_start,%d,%d,,",
		nowMillis,
		nowMillis%1000000,
	)

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"attribution_id_v2":     attribution,
			"feedback_id":           feedbackID,
			"feedback_reaction_id":  reactionID,
			"feedback_source":       "OBJECT",
			"feedback_referrer":     "/",
			"actor_id":              actorID,
			"client_mutation_id":    "1",
			"is_tracking_encrypted": true,
			"tracking":              []string{},
		},
	}

	b, err := json.Marshal(variables)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildReactionStrategies(configuredDocID string) []reactionStrategy {
	docID := strings.TrimSpace(configuredDocID)
	if docID == "" {
		docID = strings.TrimSpace(os.Getenv("FB_REACT_DOC_ID"))
	}
	strategies := make([]reactionStrategy, 0, 3)

	// If user has captured a working doc_id from browser, allow using it first.
	if docID != "" {
		strategies = append(strategies, reactionStrategy{
			name:         "comet_doc_id",
			friendlyName: cometFriendlyName,
			docID:        docID,
		})
	}

	strategies = append(strategies, reactionStrategy{
		name:         "comet_inline_query",
		friendlyName: cometFriendlyName,
		query:        cometReactionMutation,
	})
	strategies = append(strategies, reactionStrategy{
		name:         "legacy_inline_query",
		friendlyName: legacyFriendly,
		query:        legacyReactionMutation,
	})

	return strategies
}

func executeReactionRequest(cookie, fbDtsg, actorID, variablesJSON string, strategy reactionStrategy) (string, ActionErrorCode, error) {
	data := url.Values{}
	data.Set("av", actorID)
	data.Set("__user", actorID)
	data.Set("__a", "1")
	data.Set("__req", "1d")
	data.Set("fb_dtsg", fbDtsg)
	data.Set("jazoest", generateJazoest(fbDtsg))
	data.Set("fb_api_caller_class", "RelayModern")
	data.Set("fb_api_req_friendly_name", strategy.friendlyName)
	data.Set("variables", variablesJSON)
	data.Set("server_timestamps", "true")

	if strategy.docID != "" {
		data.Set("doc_id", strategy.docID)
	} else {
		data.Set("query", strategy.query)
	}

	encoded := data.Encode()
	req, err := http.NewRequest("POST", graphQLEndpoint, strings.NewReader(encoded))
	if err != nil {
		return "", ErrMalformed, fmt.Errorf("khong tao duoc request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.facebook.com/")
	req.Header.Set("Origin", "https://www.facebook.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("X-FB-Friendly-Name", strategy.friendlyName)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", ErrMalformed, fmt.Errorf("loi mang khi goi Facebook: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)
	bodyNoGuard := strings.TrimPrefix(strings.TrimSpace(body), "for (;;);")

	log.Printf("[Executor] strategy=%s status=%d body=%s", strategy.name, resp.StatusCode, TruncateStr(bodyNoGuard, 220))

	if strings.Contains(bodyNoGuard, `"error":1357032`) || strings.Contains(bodyNoGuard, `"error":1357001`) {
		return bodyNoGuard, ErrStaleContext, errors.New("fb_dtsg hoac session da het han (error 1357032 / 1357001)")
	}
	if isIncorrectQuery(bodyNoGuard) {
		return bodyNoGuard, ErrGraphqlFB, fmt.Errorf("incorrect query (%s)", strategy.name)
	}

	var parsedResp struct {
		Data *struct {
			FeedbackReact interface{} `json:"feedback_react"`
		} `json:"data"`
		Errors []interface{} `json:"errors"`
	}

	if err := json.Unmarshal([]byte(bodyNoGuard), &parsedResp); err == nil {
		if parsedResp.Data != nil && parsedResp.Data.FeedbackReact != nil {
			return bodyNoGuard, "", nil
		}
	}

	if strings.Contains(bodyNoGuard, `"errors":`) {
		return bodyNoGuard, ErrGraphqlFB, fmt.Errorf("graphql errors (%s): %s", strategy.name, TruncateStr(bodyNoGuard, 260))
	}
	if strings.Contains(bodyNoGuard, `"data":`) {
		return bodyNoGuard, "", nil
	}
	if resp.StatusCode >= 400 {
		return bodyNoGuard, ErrMalformed, fmt.Errorf("facebook tra ve status %d", resp.StatusCode)
	}

	return bodyNoGuard, ErrGraphqlFB, fmt.Errorf("phan hoi khong co data (%s): %s", strategy.name, TruncateStr(bodyNoGuard, 220))
}

func isIncorrectQuery(body string) bool {
	return strings.Contains(body, `"error":1675002`) ||
		strings.Contains(body, "Incorrect Query") ||
		strings.Contains(body, "The query provided was invalid")
}

func TruncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
