package deadmanssnitch_test

import (
	"flag"
	"math/rand"
	"testing"
	"time"

	"github.com/PremiereGlobal/go-deadmanssnitch"
)

var apiKey string
var wait int
var randomTag string
var dmsClient *deadmanssnitch.Client
var snitch deadmanssnitch.Snitch
var newSnitch *deadmanssnitch.Snitch
var updatedSnitch deadmanssnitch.Snitch
var updatedTags []string

func init() {
	var err error

	flag.StringVar(&apiKey, "apikey", "", "Dead Man's Snitch API key")
	flag.IntVar(&wait, "wait", 0, "Number of seconds to sleep before deleting snitch (so it can be manually verified)")
	flag.Parse()

	dmsClient, err = deadmanssnitch.NewClient(apiKey)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	randomTag = RandomString(10)

	snitch = deadmanssnitch.Snitch{
		Name:      "testSnitch",
		Interval:  "hourly",
		AlertType: "basic",
		Tags:      []string{"test", randomTag},
		Notes:     "This is a snitch created by github.com/PremiereGlobal/go-deadmanssnitch as a test",
	}

	updatedSnitch = deadmanssnitch.Snitch{
		Name:      "testSnitchUpdated",
		Interval:  "daily",
		AlertType: "basic",
		Tags:      []string{"testUpdated", randomTag},
		Notes:     "This is a snitch created by github.com/PremiereGlobal/go-deadmanssnitch as a test, and now it's updated",
	}

	updatedTags = []string{"newtag1", "newtag2", "newtag3"}
}

// Create a new snitch
func TestCreateSnitch(t *testing.T) {
	var err error
	newSnitch, err = dmsClient.CreateSnitch(&snitch)
	if err != nil {
		t.Error(err)
	}
}

// Check in on the snitch
func TestCheckIn(t *testing.T) {
	err := dmsClient.CheckIn(newSnitch.Token)
	if err != nil {
		t.Error(err)
	}
}

// List all of the snitches
func TestListSnitchesAll(t *testing.T) {
	snitches, err := dmsClient.ListSnitches([]string{})
	if err != nil {
		t.Error(err)
	}

	// Ensure there is at least 1 snitch (ours) and that ours is in the list
	if len(*snitches) <= 0 {
		t.Error("List all snitches - failed to find any snitches")
	}
	for _, v := range *snitches {
		if v.Token == newSnitch.Token {
			return
		}
	}

	t.Error("List all snitches - couldn't find our snitch")
}

// List just the snitches with our tag
func TestListSnitchesFiltered(t *testing.T) {
	snitches, err := dmsClient.ListSnitches([]string{"test", randomTag})
	if err != nil {
		t.Error(err)
	}

	// Ensure there is at least 1 snitch (ours) and that ours is in the list
	if len(*snitches) != 1 {
		t.Error("List filtered snitches - got more than 1 back")
	}
	for _, v := range *snitches {
		if v.Token == newSnitch.Token {
			return
		}
	}

	t.Error("List filtered snitches - couldn't find our snitch")
}

// Get the snitch we created and verify the fields
func TestGetSnitch(t *testing.T) {
	gottenSnitch, err := dmsClient.GetSnitch(newSnitch.Token)
	if err != nil {
		t.Error(err)
	}

	// Verify all the fields match what we initially created
	if gottenSnitch.Name != snitch.Name ||
		gottenSnitch.Interval != snitch.Interval ||
		gottenSnitch.AlertType != snitch.AlertType ||
		gottenSnitch.Notes != snitch.Notes ||
		!slicesEqual(gottenSnitch.Tags, snitch.Tags) {
		t.Error("Get Snitch did not match created Snitch")
	}
}

// Update the snitch we created
func TestUpdateSnitch(t *testing.T) {
	_, err := dmsClient.UpdateSnitch(newSnitch.Token, &updatedSnitch)
	if err != nil {
		t.Error(err)
	}
}

// Add tags to the snitch we created
func TestAddTags(t *testing.T) {
	err := dmsClient.AddTags(newSnitch.Token, updatedTags)
	if err != nil {
		t.Error(err)
	}
}

// Remove tags on the snitch we created
func TestRemoveTags(t *testing.T) {
	// Remove last two tags that we added
	err := dmsClient.RemoveTags(newSnitch.Token, updatedTags[len(updatedTags)-2:])
	if err != nil {
		t.Error(err)
	}
}

// Pause the snitch we created
func TestPauseSnitch(t *testing.T) {

	// Hol up
	// Seems  you can't pause a snitch right away
	time.Sleep(time.Second * 3)

	err := dmsClient.PauseSnitch(newSnitch.Token)
	if err != nil {
		t.Error(err)
	}
}

// One last check to ensure that our snitch has updated correctly and that it is paused
func TestVerifyUpdatedSnitch(t *testing.T) {
	gottenSnitch, err := dmsClient.GetSnitch(newSnitch.Token)
	if err != nil {
		t.Error(err)
	}

	// Verify all the fields match what we initially created
	expectedTags := append(updatedSnitch.Tags, updatedTags[:len(updatedTags)-2]...)
	if gottenSnitch.Name != updatedSnitch.Name ||
		gottenSnitch.Interval != updatedSnitch.Interval ||
		gottenSnitch.AlertType != updatedSnitch.AlertType ||
		gottenSnitch.Notes != updatedSnitch.Notes ||
		!slicesEqual(gottenSnitch.Tags, expectedTags) {
		t.Error("Updated Snitch did not match expected values")
	}

	// Ensure it is paused
	if gottenSnitch.Status != "paused" {
		t.Errorf("Updated Snitch is not in the paused state.  Actual state: %s", gottenSnitch.Status)
	}
}

// Delete the snitch we created
func TestDeleteSnitch(t *testing.T) {

	// Wait, if set
	time.Sleep(time.Second * time.Duration(wait))

	err := dmsClient.DeleteSnitch(newSnitch.Token)
	if err != nil {
		t.Error(err)
	}

}

func slicesEqual(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}
