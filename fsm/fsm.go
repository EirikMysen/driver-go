package fsm

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/requests"
	"fmt"
	"time"
	//må inn med request her ja.
)

//Fsm_OnItBetweenFloors
//Foreløpig får jeg ikke bruk for denne.

func ResetAndStartTimer(timer *time.Timer, duration time.Duration) {
	if !timer.Stop() {
		<-timer.C
	}
	timer.Reset(duration)
}

// Denne hører kanskje ikke helt hjemme her men.
func setAllButtonRequestsLights(elevator_M elevator.Elevator) {
	for floor := 0; floor < elevator.NumFloors; floor++ {
		for btn := 0; btn < elevator.NumButtons; btn++ {
			if elevator_M.Requests[floor][btn] == 1 {
				elevio.SetButtonLamp(elevio.ButtonType(btn), floor, true) //er btn riktig argument her da?
			} else {
				elevio.SetButtonLamp(elevio.ButtonType(btn), floor, false)

			}

		}

	}

}

func OnReQuestButtonPress(button elevio.ButtonEvent, elevator_M elevator.Elevator, timer *time.Timer, duration time.Duration) {
	print("OnRequestButtonPress is running! \n")
	switch elevator_M.Behav {
	case elevator.EB_DoorOpen:
		print("Door_Open \n")

		if requests.RequestsShouldClearImmediately(elevator_M, button.Floor) {

			ResetAndStartTimer(timer, duration)
		} else {
			elevator_M.Requests[button.Floor][button.Button] = 1
			//pair := requests.RequestsChooseDirection(elevator_M) //da blir pari automatisk av typen DirnBehaviourPair
			//elevator_M.Dirn = pair.Dirn
			//elevator_M.Behav = pair.Behaviour
			//elevio.SetMotorDirection(elevator_M.Dirn)
		}

	case elevator.EB_Moving:
		print("Door_Moving \n")

		elevator_M.Requests[button.Floor][button.Button] = 1
		//pair := requests.RequestsChooseDirection(elevator_M) //da blir pari automatisk av typen DirnBehaviourPair
		//elevator_M.Dirn = pair.Dirn
		//elevator_M.Behav = pair.Behaviour
		//elevio.SetMotorDirection(elevator_M.Dirn)

	case elevator.EB_Idle:
		print("Idle \n")

		elevator_M.Requests[button.Floor][button.Button] = 1
		pair := requests.RequestsChooseDirection(elevator_M) //da blir pari automatisk av typen DirnBehaviourPair
		elevator_M.Dirn = pair.Dirn
		elevator_M.Behav = pair.Behaviour
		switch elevator_M.Behav {

		case elevator.EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			ResetAndStartTimer(timer, duration)
			elevator_M = requests.Requests_ClearAtCurrentFloor(elevator_M)
			//heien sto stille i en etasje, så ble det trykket knapp for akkurat denne heisen. Da åpnes dørene naturligvis. trenger ikke endre retning da

		case elevator.EB_Moving:
			elevio.SetMotorDirection(elevator_M.Dirn)
			//står den stille, og så endrer RequestChooseDirection til at den skal røre på seg, må vi sette en retning ja.

		case elevator.EB_Idle:
			break
			//hvis requestChooseDirection konkluderte med fortsatt idle er dette rett.

		}

	}

	elevator.ElevatorPrint(elevator_M)
	setAllButtonRequestsLights(elevator_M)

}

func Fsm_OnFloorArrival(elevator_M elevator.Elevator, timer *time.Timer, duration time.Duration, newFloor int) {
	fmt.Printf("\nYou have arrived at floor: %d\n", newFloor)
	elevator.ElevatorPrint(elevator_M)

	elevator_M.Floor = newFloor
	elevio.SetFloorIndicator(elevator_M.Floor)

	switch elevator_M.Behav {
	case elevator.EB_Moving:
		if requests.RequestsShouldStop(elevator_M) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevator_M.Dirn = elevio.MD_Stop
			elevio.SetDoorOpenLamp(true)
			elevator_M = requests.Requests_ClearAtCurrentFloor(elevator_M)
			ResetAndStartTimer(timer, duration)
			setAllButtonRequestsLights(elevator_M)
			elevator_M.Behav = elevator.EB_DoorOpen
		}

	default:
		break

	}
	print("\n New state: \n")
	elevator.ElevatorPrint(elevator_M)

}

func Fsm_OnDoorTimeout(elevator_M elevator.Elevator, timer *time.Timer, duration time.Duration) {
	print("\n FSM_OnDoorTimeout running \n")
	elevator.ElevatorPrint(elevator_M)

	switch elevator_M.Behav {
	case elevator.EB_DoorOpen:
		pair := requests.RequestsChooseDirection(elevator_M) //da blir pari automatisk av typen DirnBehaviourPair
		elevator_M.Dirn = pair.Dirn
		elevator_M.Behav = pair.Behaviour

		switch elevator_M.Behav {
		case elevator.EB_DoorOpen:
			ResetAndStartTimer(timer, duration)
			elevator_M = requests.Requests_ClearAtCurrentFloor(elevator_M)
			setAllButtonRequestsLights(elevator_M)

		case elevator.EB_Moving, elevator.EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevator_M.Dirn)

		}
	default:
		break

	}
	print("\n New state: \n")
	elevator.ElevatorPrint(elevator_M)

}
