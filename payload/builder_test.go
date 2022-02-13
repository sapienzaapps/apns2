package payload_test

import (
	"encoding/json"
	"testing"

	. "github.com/sapienzaapps/apns2/payload"
)

func TestEmptyPayload(t *testing.T) {
	payload := NewPayload()
	b, _ := json.Marshal(payload)
	if `{"aps":{}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{}}`)
	}
}

func TestAlert(t *testing.T) {
	payload := NewPayload().Alert("hello")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":"hello"}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":"hello"}}`)
	}
}

func TestBadge(t *testing.T) {
	payload := NewPayload().Badge(1)
	b, _ := json.Marshal(payload)
	if `{"aps":{"badge":1}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"badge":1}}`)
	}
}

func TestZeroBadge(t *testing.T) {
	payload := NewPayload().ZeroBadge()
	b, _ := json.Marshal(payload)
	if `{"aps":{"badge":0}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"badge":0}}`)
	}
}

func TestUnsetBadge(t *testing.T) {
	payload := NewPayload().Badge(1).UnsetBadge()
	b, _ := json.Marshal(payload)
	if `{"aps":{}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{}}`)
	}
}

func TestSound(t *testing.T) {
	payload := NewPayload().Sound("Default.caf")
	b, _ := json.Marshal(payload)
	if `{"aps":{"sound":"Default.caf"}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"sound":"Default.caf"}}`)
	}
}

func TestSoundDictionary(t *testing.T) {
	payload := NewPayload().Sound(map[string]interface{}{
		"critical": 1,
		"name":     "default",
		"volume":   0.8,
	})
	b, _ := json.Marshal(payload)
	if `{"aps":{"sound":{"critical":1,"name":"default","volume":0.8}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"sound":{"critical":1,"name":"default","volume":0.8}}}`)
	}
}

func TestContentAvailable(t *testing.T) {
	payload := NewPayload().ContentAvailable()
	b, _ := json.Marshal(payload)
	if `{"aps":{"content-available":1}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"content-available":1}}`)
	}
}

func TestMutableContent(t *testing.T) {
	payload := NewPayload().MutableContent()
	b, _ := json.Marshal(payload)
	if `{"aps":{"mutable-content":1}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"mutable-content":1}}`)
	}
}

func TestCustom(t *testing.T) {
	payload := NewPayload().Custom("key", "val")
	b, _ := json.Marshal(payload)
	if `{"aps":{},"key":"val"}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{},"key":"val"}`)
	}
}

func TestCustomMap(t *testing.T) {
	payload := NewPayload().Custom("key", map[string]interface{}{
		"map": 1,
	})
	b, _ := json.Marshal(payload)
	if `{"aps":{},"key":{"map":1}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{},"key":{"map":1}}`)
	}
}

func TestAlertTitle(t *testing.T) {
	payload := NewPayload().AlertTitle("hello")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"title":"hello"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"title":"hello"}}}`)
	}
}

func TestAlertTitleLocKey(t *testing.T) {
	payload := NewPayload().AlertTitleLocKey("GAME_PLAY_REQUEST_FORMAT")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"title-loc-key":"GAME_PLAY_REQUEST_FORMAT"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"title-loc-key":"GAME_PLAY_REQUEST_FORMAT"}}}`)
	}
}

func TestAlertLocArgs(t *testing.T) {
	payload := NewPayload().AlertLocArgs([]string{"Jenna", "Frank"})
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"loc-args":["Jenna","Frank"]}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"loc-args":["Jenna","Frank"]}}}`)
	}
}

func TestAlertTitleLocArgs(t *testing.T) {
	payload := NewPayload().AlertTitleLocArgs([]string{"Jenna", "Frank"})
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"title-loc-args":["Jenna","Frank"]}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"title-loc-args":["Jenna","Frank"]}}}`)
	}
}

func TestAlertSubtitle(t *testing.T) {
	payload := NewPayload().AlertSubtitle("hello")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"subtitle":"hello"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"subtitle":"hello"}}}`)
	}
}

func TestAlertBody(t *testing.T) {
	payload := NewPayload().AlertBody("body")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"body":"body"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"body":"body"}}}`)
	}
}

func TestAlertLaunchImage(t *testing.T) {
	payload := NewPayload().AlertLaunchImage("Default.png")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"launch-image":"Default.png"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"launch-image":"Default.png"}}}`)
	}
}

func TestAlertLocKey(t *testing.T) {
	payload := NewPayload().AlertLocKey("LOC")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"loc-key":"LOC"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"loc-key":"LOC"}}}`)
	}
}

func TestAlertAction(t *testing.T) {
	payload := NewPayload().AlertAction("action")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"action":"action"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"action":"action"}}}`)
	}
}

func TestAlertActionLocKey(t *testing.T) {
	payload := NewPayload().AlertActionLocKey("PLAY")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"action-loc-key":"PLAY"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"action-loc-key":"PLAY"}}}`)
	}
}

func TestCategory(t *testing.T) {
	payload := NewPayload().Category("NEW_MESSAGE_CATEGORY")
	b, _ := json.Marshal(payload)
	if `{"aps":{"category":"NEW_MESSAGE_CATEGORY"}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"category":"NEW_MESSAGE_CATEGORY"}}`)
	}
}

func TestMdm(t *testing.T) {
	payload := NewPayload().Mdm("996ac527-9993-4a0a-8528-60b2b3c2f52b")
	b, _ := json.Marshal(payload)
	if `{"aps":{},"mdm":"996ac527-9993-4a0a-8528-60b2b3c2f52b"}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{},"mdm":"996ac527-9993-4a0a-8528-60b2b3c2f52b"}`)
	}
}

func TestThreadID(t *testing.T) {
	payload := NewPayload().ThreadID("THREAD_ID")
	b, _ := json.Marshal(payload)
	if `{"aps":{"thread-id":"THREAD_ID"}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"thread-id":"THREAD_ID"}}`)
	}
}

func TestURLArgs(t *testing.T) {
	payload := NewPayload().URLArgs([]string{"a", "b"})
	b, _ := json.Marshal(payload)
	if `{"aps":{"url-args":["a","b"]}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"url-args":["a","b"]}}`)
	}
}

func TestSoundName(t *testing.T) {
	payload := NewPayload().SoundName("test")
	b, _ := json.Marshal(payload)
	if `{"aps":{"sound":{"critical":1,"name":"test","volume":1}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"sound":{"critical":1,"name":"test","volume":1}}}`)
	}
}

func TestSoundVolume(t *testing.T) {
	payload := NewPayload().SoundVolume(0.5)
	b, _ := json.Marshal(payload)
	if `{"aps":{"sound":{"critical":1,"name":"default","volume":0.5}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"sound":{"critical":1,"name":"default","volume":0.5}}}`)
	}
}

func TestAlertSummaryArg(t *testing.T) {
	payload := NewPayload().AlertSummaryArg("Robert")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"summary-arg":"Robert"}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"summary-arg":"Robert"}}}`)
	}
}

func TestAlertSummaryArgCount(t *testing.T) {
	payload := NewPayload().AlertSummaryArgCount(3)
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":{"summary-arg-count":3}}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":{"summary-arg-count":3}}}`)
	}
}

func TestInterruptionLevelPassive(t *testing.T) {
	payload := NewPayload().InterruptionLevel(InterruptionLevelPassive)
	b, _ := json.Marshal(payload)
	if `{"aps":{"interruption-level":"passive"}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"interruption-level":"passive"}}`)
	}
}

func TestInterruptionLevelActive(t *testing.T) {
	payload := NewPayload().InterruptionLevel(InterruptionLevelActive)
	b, _ := json.Marshal(payload)
	if `{"aps":{"interruption-level":"active"}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"interruption-level":"active"}}`)
	}
}

func TestInterruptionLevelTimeSensitive(t *testing.T) {
	payload := NewPayload().InterruptionLevel(InterruptionLevelTimeSensitive)
	b, _ := json.Marshal(payload)
	if `{"aps":{"interruption-level":"time-sensitive"}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"interruption-level":"time-sensitive"}}`)
	}
}

func TestInterruptionLevelCritical(t *testing.T) {
	payload := NewPayload().InterruptionLevel(InterruptionLevelCritical)
	b, _ := json.Marshal(payload)
	if `{"aps":{"interruption-level":"critical"}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"interruption-level":"critical"}}`)
	}
}

func TestRelevanceScore(t *testing.T) {
	payload := NewPayload().RelevanceScore(0.1)
	b, _ := json.Marshal(payload)
	if `{"aps":{"relevance-score":0.1}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"relevance-score":0.1}}`)
	}
}

func TestRelevanceScoreZero(t *testing.T) {
	payload := NewPayload().RelevanceScore(0)
	b, _ := json.Marshal(payload)
	if `{"aps":{"relevance-score":0}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"relevance-score":0}}`)
	}
}

func TestUnsetRelevanceScore(t *testing.T) {
	payload := NewPayload().RelevanceScore(0.1).UnsetRelevanceScore()
	b, _ := json.Marshal(payload)
	if `{"aps":{}}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{}}`)
	}
}

func TestCombined(t *testing.T) {
	payload := NewPayload().Alert("hello").Badge(1).Sound("Default.caf").InterruptionLevel(InterruptionLevelActive).RelevanceScore(0.1).Custom("key", "val")
	b, _ := json.Marshal(payload)
	if `{"aps":{"alert":"hello","badge":1,"interruption-level":"active","relevance-score":0.1,"sound":"Default.caf"},"key":"val"}` != string(b) {
		t.Fatal("Expected:", string(b), " found:", 		`{"aps":{"alert":"hello","badge":1,"interruption-level":"active","relevance-score":0.1,"sound":"Default.caf"},"key":"val"}`)
	}
}
