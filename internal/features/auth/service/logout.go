package service

import "context"

func (s *AuthService) Logout(ctx context.Context, sessionId string) error {
	if err := s.cache.Delete(ctx, s.generateSessionKey(sessionId)); err != nil {
		return err
	}

	return nil
}
