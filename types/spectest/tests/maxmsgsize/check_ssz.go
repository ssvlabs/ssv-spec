package maxmsgsize

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// Return a reflect.Value created from a proper structure (and not from an interface)
func getReflectValueForObject(obj types.Encoder) reflect.Value {
	switch obj := obj.(type) {
	case *qbft.Message,
		*types.PartialSignatureMessage, *types.PartialSignatureMessages,
		*types.SSVMessage, *types.SignedSSVMessage,
		*types.ValidatorConsensusData, *types.BeaconVote:
		return reflect.ValueOf(obj).Elem()
	}
	panic("unknown type")
}

// Check if the value respects all of its ssz size tags
func checkSSZTags(t *testing.T, val reflect.Value, mustBeEqualToSSZMax bool) {
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := valType.Field(i)
		sszMaxTag := field.Tag.Get("ssz-max")
		sszSizeTag := field.Tag.Get("ssz-size")

		if sszMaxTag != "" {
			checkSSZTagForField(t, val.Field(i), field.Name, sszMaxTag, mustBeEqualToSSZMax)
		}

		if sszSizeTag != "" {
			checkSSZTagForField(t, val.Field(i), field.Name, sszSizeTag, true)
		}
	}
}

// Check if value respects the ssz size tag, if it's either a list or an array
func checkSSZTagForField(t *testing.T, fieldValue reflect.Value, fieldName, sszMaxTag string, mustBeEqual bool) {

	// Don't check fields that are not an slice or an array
	if fieldValue.Kind() != reflect.Slice && fieldValue.Kind() != reflect.Array {
		return
	}

	// Parse sizes
	sizes, err := parseSSZSizes(sszMaxTag)
	if err != nil {
		t.Errorf("Invalid ssz tag for field %s: %s", fieldName, sszMaxTag)
	}

	// Verify all lists and nested lists
	verifyQueue := initializeQueue(fieldValue)

	for len(verifyQueue) != 0 {
		next := dequeue(&verifyQueue)
		checkSize(t, next, fieldName, sizes, mustBeEqual)
		enqueueNestedLists(next, sizes, &verifyQueue)
	}
}

// Parse ssz size tags into a list of sizes.
// E.g.: "13,51852"` -> [13, 51852]
func parseSSZSizes(tag string) ([]int, error) {
	// Split the string by comma
	parts := strings.Split(tag, ",")

	// Parse each part to an integer
	var values []int
	for _, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

type sszQueueElement struct {
	Value     reflect.Value
	sizeIndex int
}

func initializeQueue(fieldValue reflect.Value) []*sszQueueElement {
	return []*sszQueueElement{{Value: fieldValue, sizeIndex: 0}}
}

func dequeue(queue *[]*sszQueueElement) *sszQueueElement {
	next := (*queue)[0]
	*queue = (*queue)[1:]
	return next
}

func checkSize(t *testing.T, elem *sszQueueElement, fieldName string, sizes []int, mustBeEqual bool) {
	if mustBeEqual {
		if elem.Value.Len() != sizes[elem.sizeIndex] {
			t.Errorf("Field %s is different than ssz tag size: %d != %d", fieldName, elem.Value.Len(), sizes[elem.sizeIndex])
		}
	} else {
		if elem.Value.Len() > sizes[elem.sizeIndex] {
			t.Errorf("Field %s is bigger than ssz tag size: %d > %d", fieldName, elem.Value.Len(), sizes[elem.sizeIndex])
		}
	}
}

func enqueueNestedLists(elem *sszQueueElement, sizes []int, queue *[]*sszQueueElement) {
	if elem.sizeIndex < len(sizes)-1 {
		nextSize := elem.sizeIndex + 1
		for nestedListIndex := 0; nestedListIndex < elem.Value.Len(); nestedListIndex++ {
			*queue = append(*queue, &sszQueueElement{Value: elem.Value.Index(nestedListIndex), sizeIndex: nextSize})
		}
	}
}
