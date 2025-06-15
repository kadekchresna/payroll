package usecase

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
	interface_attendance "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/attendance/usecase/interface"
)

type attendanceUsecase struct {
	repo interface_attendance.IAttendanceRepository
}

func NewAttendanceUsecase(repo interface_attendance.IAttendanceRepository) usecase_interface.IAttendanceUsecase {
	return &attendanceUsecase{repo: repo}
}

func (u *attendanceUsecase) CreateAttendance(ctx context.Context, a *model.Attendance) error {
	return u.repo.Create(ctx, a)
}

func (u *attendanceUsecase) GetAttendance(ctx context.Context, id int) (*model.Attendance, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *attendanceUsecase) UpdateAttendance(ctx context.Context, a *model.Attendance) error {
	return u.repo.Update(ctx, a)
}

func (u *attendanceUsecase) DeleteAttendance(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}
