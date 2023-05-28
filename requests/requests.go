package requests

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour elevator.Behaviour
}

func RequestsAbove(e elevator.Elevator) bool {
	for f := e.Floor + 1; f < elevator.NumFloors; f++ {
		for btn := 0; btn < elevator.NumButtons; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(e elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < elevator.NumButtons; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func RequestsHere(e elevator.Elevator) bool {
	for btn := 0; btn < elevator.NumButtons; btn++ {
		if e.Requests[e.Floor][btn] == 1 {
			return true
		}
	}
	return false
}

func RequestsChooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevio.MD_Up:
		if RequestsAbove(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Up,
				Behaviour: elevator.EB_Moving,
			}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Down,
				Behaviour: elevator.EB_DoorOpen,
			}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Down,
				Behaviour: elevator.EB_Moving,
			}
		} else {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Stop,
				Behaviour: elevator.EB_Idle,
			}
		}
	case elevio.MD_Down:
		if RequestsBelow(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Down,
				Behaviour: elevator.EB_Moving,
			}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Up,
				Behaviour: elevator.EB_DoorOpen,
			}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Up,
				Behaviour: elevator.EB_Moving,
			}
		} else {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Stop,
				Behaviour: elevator.EB_Idle,
			}
		}
	case elevio.MD_Stop:
		if RequestsHere(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Stop,
				Behaviour: elevator.EB_DoorOpen,
			}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Up,
				Behaviour: elevator.EB_Moving,
			}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Down,
				Behaviour: elevator.EB_Moving,
			}
		} else {
			return DirnBehaviourPair{
				Dirn:      elevio.MD_Stop,
				Behaviour: elevator.EB_Idle,
			}
		}
	default:
		return DirnBehaviourPair{
			Dirn:      elevio.MD_Stop,
			Behaviour: elevator.EB_Idle,
		}
	}
}

func RequestsShouldStop(e elevator.Elevator) bool {
	switch e.Dirn {
	case elevio.MD_Down:
		return e.Requests[e.Floor][elevio.BT_HallDown] == 1 ||
			e.Requests[e.Floor][elevio.BT_Cab] == 1 ||
			!RequestsBelow(e)
	case elevio.MD_Up:
		return e.Requests[e.Floor][elevio.BT_HallUp] == 1 ||
			e.Requests[e.Floor][elevio.BT_Cab] == 1 ||
			!RequestsAbove(e)
	default:
		return true
	}
}

func RequestsShouldClearImmediately(e elevator.Elevator, btnFloor int) bool { //dette er kun hvis man trykker på samme knapp som man er i etasjen til(?)
	return e.Floor == btnFloor
}

func Requests_ClearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	for btn := int(elevio.ButtonType(0)); btn < elevator.NumButtons; btn++ {
		e.Requests[e.Floor][btn] = 0
	}
	return e
}
