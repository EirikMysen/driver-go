package elevator

import (
	"Driver-go/driver/elevio"
	"fmt"
)

var NumFloors int = 4
var NumButtons int = 3

type Behaviour int

const (
	EB_Idle     Behaviour = 0
	EB_DoorOpen Behaviour = 1
	EB_Moving   Behaviour = 2
)

// type BehaviourChannel chan Behaviour //Om jeg må tilbake på dette

type Elevator struct {
	Floor    int
	Dirn     elevio.MotorDirection
	Requests [][]int //må initialisere med numfloor
	Behav    Behaviour
}

func InitializeElevator() Elevator {
	elevator := Elevator{
		Floor:    elevio.GetFloor(),
		Dirn:     elevio.MD_Stop,
		Requests: make([][]int, NumFloors),
		Behav:    EB_Moving,
	}

	for i := 0; i < NumFloors; i++ {
		elevator.Requests[i] = make([]int, NumButtons)
	}
	return elevator
}

func MotorDirectionToString(d elevio.MotorDirection) string {
	switch d {

	case elevio.MD_Up:
		return "MD_Up"
	case elevio.MD_Down:
		return "MD_Down"
	case elevio.MD_Stop:
		return "MD_Stop"

	default:
		return "MD UNDEFINED"
	}
}

func ButtonToString(b elevio.ButtonType) string {
	switch b {
	case elevio.BT_HallUp:
		return "BT_HallUp"
	case elevio.BT_HallDown:
		return "BT_HallDown"
	case elevio.BT_Cab:
		return "B_Cab"
	default:
		return "B UNDEFINED"
	}
}

func BehaviourtoString(b Behaviour) string {
	switch b {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB UNDEFINED"
	}
}

func ElevatorPrint(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |floor = %-2d          |\n", es.Floor)
	fmt.Printf("  |dirn  = %-12.12s|\n", MotorDirectionToString(es.Dirn))
	fmt.Printf("  |behav = %-12.12s|\n", BehaviourtoString(es.Behav))
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := NumFloors - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < NumButtons; btn++ {
			if (f == NumFloors-1 && btn == int(elevio.BT_HallUp)) ||
				(f == 0 && btn == elevio.BT_HallDown) {
				fmt.Print("|     ")
			} else {
				fmt.Print(func() string {
					if es.Requests[f][btn] != 0 {
						return "|  #  "
					} else {
						return "|  -  "
					}
				}())
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}
