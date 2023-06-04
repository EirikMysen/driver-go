package main

import (
	"Driver-go/driver/elevator"
	"Driver-go/driver/elevio"
	"Driver-go/driver/fsm"
	"Driver-go/driver/requests"
	"Driver-go/network/bcast"
	"Driver-go/network/localip"
	"Driver-go/network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)

//HEIS 1/AVSENDER

func main() {

	elevio.Init("localhost:15657", elevator.NumFloors)
	inputPollRateMs := 25

	var d elevio.MotorDirection = elevio.MD_Down
	elevio.SetMotorDirection(d)
	e := elevator.InitializeElevator()
	ObstructionSwitch := false
	//initaliseres med getFloor, retning er MD_Stop, og Behav er EB_MOVING. Får se om dette byr på problemer..
	duration := 3 * time.Second
	Timer := time.NewTimer(duration)
	Network_Active := false
	//Network part

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	peerUpdateCh := make(chan peers.PeerUpdate)

	peerTxEnable := make(chan bool)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	helloTx := make(chan elevator.Elevator)
	helloRx := make(chan elevator.Elevator)

	go bcast.Transmitter(16570, helloTx)
	go bcast.Receiver(16570, helloRx)

	go func() {
		elevator_network := e

		for {

			helloTx <- elevator_network
			time.Sleep(1 * time.Second)
		}
	}()

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

		select {
		case a := <-drv_buttons:
			//drvButtonsSelected = true ...
			fmt.Printf("%+v\n", a)
			fmt.Printf("%+v\n", "buttons\n")
			elevio.SetButtonLamp(a.Button, a.Floor, true)

			if !ObstructionSwitch {

				e = fsm.OnReQuestButtonPress(a, e, Timer, duration)
			} else {
				e.Requests[a.Floor][a.Button] = 1
				//mulig det blir rør om du trykker på den etasjen du er i. får se
			}

			//For denne casen, så må det at en knapp blir trykket inn, påvirke  retning, state til Elevator, ButtonType(i allerede)
			//Spørmålet er ,her eller i PollButtons?

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a) //mtp å kvitte seg med ordre/registrere at ordren er fullført, mangelfullt?
			fmt.Printf("%+v\n", "floor")

			e = fsm.Fsm_OnFloorArrival(e, Timer, duration, a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			fmt.Printf("%+v\n", "obstructions")
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
				ObstructionSwitch = true
			} else {

				ObstructionSwitch = false
				pair := requests.RequestsChooseDirection(e) //da blir pari automatisk av typen DirnBehaviourPair
				e.Dirn = pair.Dirn
				e.Behav = pair.Behaviour
				elevio.SetMotorDirection(e.Dirn)

			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			fmt.Printf("%+v\n", "stop")
			for f := 0; f < elevator.NumFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
					e.Requests[f][b] = 0 //sletter alle ordre.

				}

			}
			elevio.SetMotorDirection(elevio.MD_Stop)
			e.Behav = elevator.EB_Idle
		default:
		}

		select {
		case <-Timer.C:
			if !ObstructionSwitch {
				print()
				fmt.Println("Timer is done")
				Timer.Stop()
				//Timer.Reset(duration)
				e = fsm.Fsm_OnDoorTimeout(e, Timer, duration)

			}

		default:
		}
		if Network_Active {
			select {
			case p := <-peerUpdateCh:
				fmt.Printf("Peer update:\n")
				fmt.Printf("  Peers:    %q\n", p.Peers)
				fmt.Printf("  New:      %q\n", p.New)
				fmt.Printf("  Lost:     %q\n", p.Lost)

			case a := <-helloRx:
				fmt.Printf("Received: %#v\n", a)
			}

			time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
		}
	}
}
