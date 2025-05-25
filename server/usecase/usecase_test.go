package usecase

import (
	"fmt"
	"testing"
	"time"

	"hmb_fighting/server/db"
	"hmb_fighting/server/entities"
	"hmb_fighting/server/types"
)

func TestLoginRejectsWrongPassword(t *testing.T) {
	uc := NewUsecase(db.NewMockDatabase())
	email := fmt.Sprintf("login-%d@example.com", time.Now().UnixNano())
	if _, err := uc.RegisterUser(entities.User{Name: "User", Email: email, Password: "correct-password"}); err != nil {
		t.Fatal(err)
	}

	if _, err := uc.LoginUser(email, "wrong-password"); err == nil {
		t.Fatal("expected wrong password to be rejected")
	}
}

func TestRestartAndLeaveRoomReturn(t *testing.T) {
	uc := NewUsecase(db.NewMockDatabase())
	email := fmt.Sprintf("room-%d@example.com", time.Now().UnixNano())
	user, err := uc.RegisterUser(entities.User{Name: "User", Email: email, Password: "correct-password"})
	if err != nil {
		t.Fatal(err)
	}
	roomID, err := uc.CreateRoom(user.AccessToken)
	if err != nil {
		t.Fatal(err)
	}

	assertReturns(t, func() error { return uc.RestartRoom(user.AccessToken, roomID) })
	assertReturns(t, func() error { return uc.LeaveRoom(user.AccessToken, roomID) })
}

func TestAttackIgnoresCharacterOutsideBoard(t *testing.T) {
	attacker := entities.Character{ID: 1, TeamID: 0, HP: 100, Position: [2]int{0, 0}, Weapon: "sword"}
	target := entities.Character{ID: 2, TeamID: 1, HP: 100, Position: [2]int{-1, -1}}
	room := &entities.Room{
		Phase:         types.GamePhaseAction,
		Teams:         map[int]entities.Team{0: {Characters: []entities.Character{attacker}}, 1: {Characters: []entities.Character{target}}},
		WeaponsConfig: map[string]entities.Weapon{"sword": {Range: 1}},
	}

	NewUsecase(db.NewMockDatabase()).handleAttackAction(room, room.FindCharacter(1), entities.Action{TargetID: 2})

	if got := room.FindCharacter(2).HP; got != 100 {
		t.Fatalf("got target HP %d, want 100", got)
	}
}

func assertReturns(t *testing.T, fn func() error) {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("operation did not return")
	}
}
