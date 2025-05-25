package types

type BoardTypes [BoardVerticalSize][BoardHorizontalSize]int // Обновляем с [20][10] на [16][9]

const BoardVerticalSize = 16
const BoardHorizontalSize = 9

type PhaseTypes string

type OpportunityAttackTypes string

type UserRoleTypes string

type ActionTypes string

const (
	ActionPlace   ActionTypes = "place"
	ActionStart   ActionTypes = "start"
	ActionMove    ActionTypes = "move"
	ActionAttack  ActionTypes = "attack"
	ActionAbility ActionTypes = "ability"
	ActionEndTurn ActionTypes = "end_turn"
)
const (
	UserRoleSpectator UserRoleTypes = "spectator"
	UserRolePlayer    UserRoleTypes = "player"
)

const (
	OpportunityAttackTrip   OpportunityAttackTypes = "trip"
	OpportunityAttackAttack OpportunityAttackTypes = "attack"
)

const (
	GamePhaseFinished PhaseTypes = "finished"
	GamePhasePickTeam PhaseTypes = "pick_team"
	GamePhaseMove     PhaseTypes = "move"
	GamePhaseSetup    PhaseTypes = "setup"
	GamePhaseAction   PhaseTypes = "action"
)
