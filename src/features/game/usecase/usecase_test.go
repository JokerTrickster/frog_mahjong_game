package usecase

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedDora 함수를 테스트코드를 작성합니다.
// 매개변수 request.ScoreCard 와 mysql.Cards 를 받아서 dora 카드가 있는지 확인하는 함수입니다.
// 테스트 케이스:
// - dora 카드가 있는 경우
// - dora 카드가 없는 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

import (
	"main/features/game/model/request"
	"main/utils/db/mysql"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestIsCheckedDora(t *testing.T) {
	tests := []struct {
		name string
		req  request.ScoreCard
		card mysql.Cards
		want bool
	}{
		{
			name: "dora 카드가 있는 경우",
			req: request.ScoreCard{
				CardID: 1,
				State:  "none",
				Name:   "1",
				Color:  "red",
			},
			card: mysql.Cards{
				Name:  "1",
				Color: "red",
				State: "none",
			},
			want: true,
		},
		{
			name: "dora 카드가 없는 경우",
			req: request.ScoreCard{
				CardID: 1,
				State:  "none",
				Name:   "1",
				Color:  "red",
			},
			card: mysql.Cards{
				Name:  "2",
				Color: "red",
				State: "none",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedDora(tt.req, tt.card)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedChanTa 함수를 테스트코드를 작성합니다.
// 매개변수 길이가 6인 request.ScoreCard 슬라이스를 받아서 찬타 조건을 만족하는지 확인하는 함수입니다.
// 테스트 케이스:
// - 찬타 조건을 만족하는 경우
// - 찬타 조건을 만족하지 않는 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedChanTa(t *testing.T) {
	tests := []struct {
		name  string
		cards []request.ScoreCard
		want  bool
	}{
		{
			name: "찬타 조건을 만족하는 경우",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "2",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "3",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "7",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "8",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "9",
					Color:  "red",
				},
			},
			want: true,
		},
		{
			name: "찬타 조건을 만족하지 않는 경우",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "2",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "3",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "4",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "5",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "6",
					Color:  "red",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedChanTa(tt.cards)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedTangYaoCard 함수를 테스트코드를 작성합니다.
// 매개변수 길이가 6인 request.ScoreCard 슬라이스를 받아서 탕야오 조건을 만족하는지 확인하는 함수입니다.
// 테스트 케이스:
// - 탕야오 조건을 만족하는 경우
// - 탕야오 조건을 만족하지 않는 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedTangYaoCard(t *testing.T) {
	tests := []struct {
		name  string
		cards []request.ScoreCard
		want  bool
	}{
		{
			name: "탕야오 조건을 만족하는 경우",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "2",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "3",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "4",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "5",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "6",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "7",
					Color:  "red",
				},
			},
			want: true,
		},
		{
			name: "탕야오 조건을 만족하지 않는 경우",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "3",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "4",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "5",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "6",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "7",
					Color:  "red",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedTangYaoCard(tt.cards)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedRedCard 함수를 테스트코드를 작성합니다.
// 매개변수 request.ScoreCard 를 받아서 적패 조건을 만족하는지 확인하는 함수입니다.
// 테스트 케이스:
// - 적패 조건을 만족하는 경우
// - 적패 조건을 만족하지 않는 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedRedCard(t *testing.T) {
	tests := []struct {
		name string
		card request.ScoreCard
		want bool
	}{
		{
			name: "적패 조건을 만족하는 경우",
			card: request.ScoreCard{
				CardID: 1,
				State:  "none",
				Name:   "1",
				Color:  "red",
			},
			want: true,
		},
		{
			name: "적패 조건을 만족하지 않는 경우",
			card: request.ScoreCard{
				CardID: 1,
				State:  "none",
				Name:   "1",
				Color:  "green",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedRedCard(tt.card)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedNumberCard 함수를 테스트코드를 작성합니다.
// 매개변수 길이가 6인 request.ScoreCard 를 받아서 1~9까지 숫자로만 이루어져 있는지 확인하는 함수입니다.
// 테스트 케이스:
// - 1~9까지 숫자로만 이루어져 있는 경우
// - 1~9까지 숫자로만 이루어져 있지 않은 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedNumberCard(t *testing.T) {
	tests := []struct {
		name  string
		cards []request.ScoreCard
		want  bool
	}{
		{
			name: "1~9까지 숫자로만 이루어져 있는 경우",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "2",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "3",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "4",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "5",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "6",
					Color:  "red",
				},
			},
			want: true,
		},
		{
			name: "1~9까지 숫자로만 이루어져 있지 않은 경우",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "2",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "3",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "4",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "5",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "중",
					Color:  "red",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedNumberCard(tt.cards)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedChinYao 함수를 테스트코드를 작성합니다.
// 매개변수 길이가 6인 request.ScoreCard 를 받아서 모든 패가 1,9,중,발로만 이루어진 카드인지 확인하는 함수입니다.
// 테스트 케이스:
// - 모든 패가 1,9,중,발로만 이루어진 카드인 경우
// - 모든 패가 1,9,중,발로만 이루어진 카드가 아닌 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedChinYao(t *testing.T) {
	tests := []struct {
		name  string
		cards []request.ScoreCard
		want  bool
	}{
		{
			name: "모든 패가 1,9,중,발로만 이루어진 카드인 경우",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "9",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "중",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "발",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "9",
					Color:  "red",
				},
			},
			want: true,
		},
		{
			name: "모든 패가 1,9,중,발로만 이루어진 카드가 아닌 경우",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "9",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "중",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "발",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "2",
					Color:  "red",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedChinYao(tt.cards)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedAllRed 함수를 테스트코드를 작성합니다.
// 매개변수 길이가 6인 request.ScoreCard 를 받아서 모든 패가 red이면 true 아니라면 false를 반환하는 함수입니다.
// 테스트 케이스:
// - 모든 패가 red이면 true
// - 모든 패가 red가 아니면 false
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedAllRed(t *testing.T) {
	tests := []struct {
		name  string
		cards []request.ScoreCard
		want  bool
	}{
		{
			name: "모든 패가 red이면 true",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "2",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "3",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "4",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "5",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "6",
					Color:  "red",
				},
			},
			want: true,
		},
		{
			name: "모든 패가 red가 아니면 false",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "red",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "2",
					Color:  "red",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "3",
					Color:  "red",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "4",
					Color:  "red",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "5",
					Color:  "red",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "6",
					Color:  "green",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedAllRed(tt.cards)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedAllGreen 함수를 테스트코드를 작성합니다.
// 매개변수 길이가 6인 request.ScoreCard 를 받아서 모든 패가 green이면 true 아니라면 false를 반환하는 함수입니다.
// 테스트 케이스:
// - 모든 패가 green이면 true
// - 모든 패가 green이 아니면 false
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedAllGreen(t *testing.T) {
	tests := []struct {
		name  string
		cards []request.ScoreCard
		want  bool
	}{
		{
			name: "모든 패가 green이면 true",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "green",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "2",
					Color:  "green",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "3",
					Color:  "green",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "4",
					Color:  "green",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "5",
					Color:  "green",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "6",
					Color:  "green",
				},
			},
			want: true,
		},
		{
			name: "모든 패가 green이 아니면 false",
			cards: []request.ScoreCard{
				{
					CardID: 1,
					State:  "none",
					Name:   "1",
					Color:  "green",
				},
				{
					CardID: 2,
					State:  "none",
					Name:   "2",
					Color:  "green",
				},
				{
					CardID: 3,
					State:  "none",
					Name:   "3",
					Color:  "green",
				},
				{
					CardID: 4,
					State:  "none",
					Name:   "4",
					Color:  "green",
				},
				{
					CardID: 5,
					State:  "none",
					Name:   "5",
					Color:  "green",
				},
				{
					CardID: 6,
					State:  "none",
					Name:   "6",
					Color:  "red",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedAllGreen(tt.cards)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedContinuousCard 함수를 테스트코드를 작성합니다.
// 매개변수 request.ScoreCard 타입 3개를 받아서 3개 카드가 연속된 숫자라면 true 아니라면 false를 반환하는 함수입니다.
// 테스트 케이스:
// - 3개 카드가 연속된 숫자인 경우
// - 3개 카드가 연속된 숫자가 아닌 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedContinuousCard(t *testing.T) {
	tests := []struct {
		name  string
		card1 request.ScoreCard
		card2 request.ScoreCard
		card3 request.ScoreCard
		want  bool
	}{
		{
			name: "3개 카드가 연속된 숫자인 경우",
			card1: request.ScoreCard{
				CardID: 1,
				State:  "none",
				Name:   "1",
				Color:  "red",
			},
			card2: request.ScoreCard{
				CardID: 2,
				State:  "none",
				Name:   "2",
				Color:  "red",
			},
			card3: request.ScoreCard{
				CardID: 3,
				State:  "none",
				Name:   "3",
				Color:  "red",
			},
			want: true,
		},
		{
			name: "3개 카드가 연속된 숫자가 아닌 경우",
			card1: request.ScoreCard{
				CardID: 1,
				State:  "none",
				Name:   "1",
				Color:  "red",
			},
			card2: request.ScoreCard{
				CardID: 2,
				State:  "none",
				Name:   "2",
				Color:  "red",
			},
			card3: request.ScoreCard{
				CardID: 3,
				State:  "none",
				Name:   "4",
				Color:  "red",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedContinuousCard(tt.card1, tt.card2, tt.card3)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedSameCard 함수를 테스트코드를 작성합니다.
// 매개변수 request.ScoreCard 타입 3개를 받아서 3개 카드가 모두 같은 카드라면 true 아니라면 false를 반환하는 함수입니다.
// 테스트 케이스:
// - 모두 같은 카드인 경우
// - 모두 같은 카드가 아닌 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedSameCard(t *testing.T) {
	tests := []struct {
		name  string
		card1 request.ScoreCard
		card2 request.ScoreCard
		card3 request.ScoreCard
		want  bool
	}{
		{
			name: "모두 같은 카드인 경우",
			card1: request.ScoreCard{
				CardID: 1,
				State:  "none",
				Name:   "1",
				Color:  "red",
			},
			card2: request.ScoreCard{
				CardID: 2,
				State:  "none",
				Name:   "1",
				Color:  "normal",
			},
			card3: request.ScoreCard{
				CardID: 3,
				State:  "none",
				Name:   "1",
				Color:  "normal",
			},
			want: true,
		},
		{
			name: "모두 같은 카드가 아닌 경우",
			card1: request.ScoreCard{
				CardID: 1,
				State:  "none",
				Name:   "1",
				Color:  "red",
			},
			card2: request.ScoreCard{
				CardID: 2,
				State:  "none",
				Name:   "2",
				Color:  "red",
			},
			card3: request.ScoreCard{
				CardID: 3,
				State:  "none",
				Name:   "3",
				Color:  "red",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedSameCard(tt.card1, tt.card2, tt.card3)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}

// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// when-then 패턴을 사용하여 테스트를 작성합니다.
// src/features/game/usecase/usecase.go 파일의 IsCheckedWinRequest 함수를 테스트코드를 작성합니다.
// 매개변수 mysql.RoomUsers, score int 를 받아서 승리 요청이 가능한지 확인하는 함수입니다.
// 테스트 케이스:
// - 플레이 상태가 play이고 소유 카드가 6장이며 점수가 10점인 경우
// - 플레이 상태가 loan이고 소유 카드가 6장이며 점수가 5점인 경우
// - 플레이 상태가 play-wait이고 소유 카드가 5장이며 점수가 5점인 경우
// - 플레이 상태가 play-wait이고 소유 카드가 5장이며 점수가 0점인 경우
// 테스트 경로: src/features/game/usecase/usecase_test.go

func TestIsCheckedWinRequest(t *testing.T) {

	tests := []struct {
		name     string
		roomUser mysql.RoomUsers
		score    int
		want     bool
	}{
		{
			name: "플레이 상태가 play이고 소유 카드가 6장이며 점수가 10점인 경우",
			roomUser: mysql.RoomUsers{
				PlayerState:    "play",
				OwnedCardCount: 6,
			},
			score: 10,
			want:  true,
		},
		{
			name: "플레이 상태가 loan이고 소유 카드가 6장이며 점수가 5점인 경우",
			roomUser: mysql.RoomUsers{
				PlayerState:    "loan",
				OwnedCardCount: 6,
			},
			score: 5,
			want:  true,
		},
		{
			name: "플레이 상태가 play-wait이고 소유 카드가 5장이며 점수가 5점인 경우",
			roomUser: mysql.RoomUsers{
				PlayerState:    "play-wait",
				OwnedCardCount: 5,
			},
			score: 5,
			want:  false,
		},
		{
			name: "플레이 상태가 play-wait이고 소유 카드가 5장이며 점수가 0점인 경우",
			roomUser: mysql.RoomUsers{
				PlayerState:    "play-wait",
				OwnedCardCount: 5,
			},
			score: 0,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := IsCheckedWinRequest(tt.roomUser, tt.score)

			// then
			assert.Equal(t, got, tt.want)
		})
	}
}
