package database

import (
	"fmt"

	supabase "github.com/supabase-community/supabase-go"
)

type SupabaseStore struct {
	Client *supabase.Client
}

func NewSupabase(url, key string) (*SupabaseStore, error) {
	if url == "" || key == "" {
		return nil, fmt.Errorf("Supabase URL 또는 Key가 설정되지 않았습니다")
	}

	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		return nil, fmt.Errorf("supabase 클라이언트 생성 실패: %w", err)
	}

	return &SupabaseStore{Client: client}, nil
}
