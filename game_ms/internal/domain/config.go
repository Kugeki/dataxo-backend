package domain

type DisappearingModeConfig struct {
	// 0 is no limit
	PlayerFiguresLimit int
	WinLineLength      int
	BoardWidth         int
	BoardHeight        int
}

func ValidateDisappearingModeConfig(cfg DisappearingModeConfig) error {
	if cfg.PlayerFiguresLimit < 0 {
		return ErrNegativePlayerFiguresLimit
	}
	if cfg.WinLineLength <= 0 {
		return ErrNegativeOrZeroedWinLineLength
	}
	if cfg.BoardWidth <= 0 {
		return ErrNegativeOrZeroedBoardWidth
	}
	if cfg.BoardHeight <= 0 {
		return ErrNegativeOrZeroedBoardHeight
	}
	return nil
}
