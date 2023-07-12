package service

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"log"
	"net"
	"time"
)

type Service struct {
	tvEndpoint   string
	tvAuth       string
	wolBroadcast string
	wolMac       net.HardwareAddr

	wolClient WolClient
	tvClient  TVClient
}

type WolClient interface {
	Wake(addr string, target net.HardwareAddr) error
}

type TVClient interface {
	KeyExit() error
	KeyExternalInput() error
	KeyRight() error
	KeyLeft() error
	KeyOK() error
	GetScreen() (image.Image, error)
}

func New(wc WolClient, tv TVClient, tvEndpoint, tvAuth, wolBroadcast, wolMac string) (*Service, error) {
	mac, err := net.ParseMAC(wolMac)
	if err != nil {
		return nil, err
	}

	return &Service{
		tvEndpoint:   tvEndpoint,
		tvAuth:       tvAuth,
		wolBroadcast: wolBroadcast,
		wolMac:       mac,

		wolClient: wc,
		tvClient:  tv,
	}, nil
}

func (svc *Service) Wake() error {
	if err := svc.wolClient.Wake(svc.wolBroadcast, svc.wolMac); err != nil {
		return err
	}

	return nil
}

func (svc *Service) IsPoweredOn() bool {
	t := time.Duration(2) * time.Second
	_, err := net.DialTimeout("tcp", svc.tvEndpoint, t)

	return err == nil
}

func drawSlot(dst draw.Image, slot int) {
	x := 65 + slot*140
	draw.Draw(
		dst,
		image.Rect(x, dst.Bounds().Max.Y-225, x+5, dst.Bounds().Max.Y-220),
		&image.Uniform{color.Black},
		image.ZP,
		draw.Over,
	)
}

func isActive(img image.Image, slot int) bool {
	x := 65 + slot*140
	y := img.Bounds().Max.Y - 225

	col := img.At(x, y)

	r, g, b, _ := col.RGBA()

	redOK := r < 63572
	greenOK := g < 58935
	blueOK := b < 52623

	log.Printf("SLOT %d red: %d, green: %d, blue: %d, OK: %v\n", slot, r, g, b, redOK && greenOK && blueOK)
	return redOK && greenOK && blueOK
}

func (s *Service) getActiveSlot() (int, error) {
	img, err := s.tvClient.GetScreen()
	if err != nil {
		return 0, err
	}

	for i := 0; i < 6; i++ {
		if isActive(img, i) {
			return i, nil
		}
	}

	return 0, errors.New("no active slot")
}

func (s *Service) SwitchToPS5() error {
	s.tvClient.KeyExit()
	time.Sleep(500 * time.Millisecond)
	s.tvClient.KeyExternalInput()
	time.Sleep(3 * time.Second)

	slot, err := s.getActiveSlot()
	if err != nil {
		panic(err)
	}

	if slot < 2 {
		for x := slot; x < 2; x++ {
			s.tvClient.KeyRight()
			time.Sleep(500 * time.Millisecond)
		}
	} else if slot > 2 {
		for x := slot; x > 1; x-- {
			s.tvClient.KeyLeft()
			time.Sleep(500 * time.Millisecond)
		}
	}

	return s.tvClient.KeyOK()
}

func (s *Service) SwitchToFireTV() error {
	s.tvClient.KeyExit()
	time.Sleep(500 * time.Millisecond)
	s.tvClient.KeyExternalInput()
	time.Sleep(3 * time.Second)

	slot, err := s.getActiveSlot()
	if err != nil {
		panic(err)
	}

	if slot != 1 {
		if slot == 0 {
			s.tvClient.KeyRight()
			time.Sleep(500 * time.Millisecond)
		}

		for x := slot; x > 1; x-- {
			s.tvClient.KeyLeft()
			time.Sleep(500 * time.Millisecond)
		}
	}

	return s.tvClient.KeyOK()
}

func (s *Service) SwitchToTV() error {
	s.tvClient.KeyExit()
	time.Sleep(500 * time.Millisecond)
	s.tvClient.KeyExternalInput()
	time.Sleep(3 * time.Second)

	slot, err := s.getActiveSlot()
	if err != nil {
		panic(err)
	}

	if slot != 0 {
		for x := slot; x > 0; x-- {
			s.tvClient.KeyLeft()
			time.Sleep(500 * time.Millisecond)
		}
	}

	return s.tvClient.KeyOK()
}
