package usecases

import (
	"context"
	"fmt"
	"sort"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

type SessionUseCase struct {
	trainerRepo repositories.TrainerRepository
}

func ProvideSessionUseCase(trainerRepo repositories.TrainerRepository) *SessionUseCase {
	return &SessionUseCase{
		trainerRepo: trainerRepo,
	}
}

// MatchTrainer - จับคู่เทรนเนอร์ที่เหมาะสมตาม Use Case 4S
// Algorithm:
// 1. หาเทรนเนอร์ที่ว่างในวันและเวลาที่กำหนด
// 2. นับจำนวนนัดหมายของแต่ละเทรนเนอร์ในวันนั้น
// 3. เรียงตามจำนวนนัดหมาย (น้อย -> มาก) และวันที่สร้าง (เก่า -> ใหม่)
// 4. ตรวจสอบว่ามีนัดซ้อนทับหรือไม่
// 5. คืนค่าเทรนเนอร์คนแรกที่ไม่มีนัดซ้อนทับ
func (u *SessionUseCase) MatchTrainer(ctx context.Context, req requests.MatchTrainerRequest) (*responses.TrainerMatchResponse, error) {
	// Step 1: Find available trainers by day and time
	availableTrainers, err := u.trainerRepo.FindAvailableTrainers(ctx, req.DayOfWeek, req.StartTime, req.EndTime)
	fmt.Println("availableTrainers", availableTrainers)
	if err != nil {
		return nil, fmt.Errorf("failed to find available trainers: %w", err)
	}

	if len(availableTrainers) == 0 {
		return nil, fmt.Errorf("NO_TRAINER_AVAILABLE")
	}

	// Step 2 & 3: Count appointments and sort
	type trainerWithCount struct {
		info  repositories.TrainerInfo
		count int64
	}

	trainersWithCounts := make([]trainerWithCount, 0, len(availableTrainers))
	for _, trainer := range availableTrainers {
		count, err := u.trainerRepo.CountAppointmentsOnDate(ctx, trainer.Username, req.StartTime)
		if err != nil {
			// If error, assume 0 appointments
			count = 0
		}

		trainersWithCounts = append(trainersWithCounts, trainerWithCount{
			info:  trainer,
			count: count,
		})
	}
	fmt.Println("trainersWithCounts", trainersWithCounts)

	// Sort by appointment count (ASC), then by created_at (ASC)
	sort.Slice(trainersWithCounts, func(i, j int) bool {
		if trainersWithCounts[i].count != trainersWithCounts[j].count {
			return trainersWithCounts[i].count < trainersWithCounts[j].count
		}
		return trainersWithCounts[i].info.CreatedAt.Before(trainersWithCounts[j].info.CreatedAt)
	})
	fmt.Println("trainersWithCounts sorted", trainersWithCounts)

	// Step 4 & 5: Check for overlaps and return first available
	for _, tc := range trainersWithCounts {
		hasOverlap, err := u.trainerRepo.CheckScheduleOverlap(ctx, tc.info.Username, req.StartTime, req.EndTime)
		fmt.Println("hasOverlap", hasOverlap)
		fmt.Println("err", err)
		if err != nil {
			// If error checking overlap, skip this trainer
			continue
		}

		if !hasOverlap {
			// Found a trainer without overlap!
			fmt.Println("tc", tc)
			return &responses.TrainerMatchResponse{
				TrainerUsername: tc.info.Username,
				TrainerName:     fmt.Sprintf("%s %s", tc.info.FirstName, tc.info.LastName),
				DayOfWeek:       req.DayOfWeek,
				StartTime:       req.StartTime,
				EndTime:         req.EndTime,
				Appointments:    tc.count,
			}, nil
		}
	}

	// No trainer found without overlap
	return nil, fmt.Errorf("NO_TRAINER_AVAILABLE")
}

// ListAllTrainers - ดึงรายชื่อเทรนเนอร์ทั้งหมด (สำหรับ dropdown)
func (u *SessionUseCase) ListAllTrainers(ctx context.Context) ([]responses.TrainerListResponse, error) {
	trainers, err := u.trainerRepo.ListAllTrainers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list trainers: %w", err)
	}

	result := make([]responses.TrainerListResponse, len(trainers))
	for i, t := range trainers {
		result[i] = responses.TrainerListResponse{
			Username: t.Username,
			Name:     fmt.Sprintf("%s %s", t.FirstName, t.LastName),
		}
	}

	return result, nil
}
