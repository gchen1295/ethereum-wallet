package piratesea

import "time"

func getCurrentTS() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000")
}

func newPollNewListingsPayload(collections []string, count int, lastPolled string) *EventHistoryQueryVariables {
	if count < 100 {
		count = 100
	}

	return &EventHistoryQueryVariables{
		Collections:    collections,
		Count:          count,
		EventTimestamp: lastPolled,
		ShowAll:        true,
		EventTypes:     []string{"AUCTION_CREATED"},
	}
}
