package withfriend

import "errors"

type Config struct {
	// 0 is no limit
	PlayerFiguresLimit int
	WinLineLength      int
	BoardWidth         int
	BoardHeight        int
}

func ValidateConfig(cfg Config) error {
	if cfg.PlayerFiguresLimit < 0 {
		return errors.New("player figures limit is negative")
	}
	if cfg.WinLineLength <= 0 {
		return errors.New("win line length is negative or equals to zero")
	}
	if cfg.BoardWidth <= 0 {
		return errors.New("board width is negative or equals to zero")
	}
	if cfg.BoardHeight <= 0 {
		return errors.New("board height is negative or equals to zero")
	}
	return nil
}
