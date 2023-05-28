package main

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"time"
)

func main() {

	elevio.Init("localhost:15657", elevator.NumFloors)

	var d elevio.MotorDirection = elevio.MD_Down
	elevio.SetMotorDirection(d)
	e := elevator.InitializeElevator()
	//initaliseres med getFloor, retning er MD_Down, og Behav er EB_MOVING. Får se om dette byr på problemer..
	duration := 3 * time.Second
	Timer := time.NewTimer(duration)

	//drvButtonsSelected := false
	//Her i koden har main.c en funksjon som sjekker om heisen er imellom heiser. endrer da retning til ned ,og behaviour til Moving. Har satt dette
	//allerede, får se om det skaper problemer.
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons) //endrer drv_buttons hvis en knapp ble trøkket inn. av typen ButtonEvent(floor,buttonType)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		print("Cases button/floor/stop/obstc STARTER \n")
		select {
		case a := <-drv_buttons:
			//drvButtonsSelected = true ...
			fmt.Printf("%+v\n", a)
			fmt.Printf("%+v\n", "buttons\n")
			elevio.SetButtonLamp(a.Button, a.Floor, true)

			fsm.OnReQuestButtonPress(a, e, Timer, duration)
			//For denne casen, så må det at en knapp blir trykket inn, påvirke  retning, state til Elevator, ButtonType(i allerede)
			//Spørmålet er ,her eller i PollButtons?

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a) //mtp å kvitte seg med ordre/registrere at ordren er fullført, mangelfullt?
			fmt.Printf("%+v\n", "floor")

			fsm.Fsm_OnFloorArrival(e, Timer, duration, a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			fmt.Printf("%+v\n", "obstructions")
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			fmt.Printf("%+v\n", "stop")
			for f := 0; f < elevator.NumFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
		e.Floor = elevio.GetFloor()
		print("Cases button/floor/stop/obstc FERDIG \n")
		select {
		case <-Timer.C:
			print()
			fmt.Println("Timer is done")
			Timer.Stop()
			fsm.Fsm_OnDoorTimeout(e, Timer, duration)
		default:

		}

	}
}
