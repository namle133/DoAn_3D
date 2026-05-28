package services

import (
	"log"
	"time"
)

// SchedulerService runs periodic background tasks
type SchedulerService struct {
	stopChannel chan bool
	isRunning   bool
}

var scheduler *SchedulerService

func GetScheduler() *SchedulerService {
	if scheduler == nil {
		scheduler = &SchedulerService{
			stopChannel: make(chan bool),
			isRunning:   false,
		}
	}
	return scheduler
}

// StartScheduler starts all background jobs
func (s *SchedulerService) StartScheduler() {
	if s.isRunning {
		return
	}

	s.isRunning = true
	log.Println("✓ Scheduler started")

	// Run daily invoice generation at 00:00
	go s.scheduleDaily(0, func() {
		log.Println("Running: Monthly invoice generation")
		invoiceService := NewInvoiceService()
		currentMonth := time.Now().Format("2006-01")
		if err := invoiceService.GenerateMonthlyInvoices(currentMonth); err != nil {
			log.Printf("Error generating invoices: %v", err)
		}
	})

	// Run daily overdue invoice check at 08:00
	go s.scheduleDaily(8*time.Hour, func() {
		log.Println("Running: Overdue invoice check")
		invoiceService := NewInvoiceService()
		if err := invoiceService.CheckOverdueInvoices(); err != nil {
			log.Printf("Error checking overdue invoices: %v", err)
		}
	})

	// Run daily contract expiry check at 09:00
	go s.scheduleDaily(9*time.Hour, func() {
		log.Println("Running: Contract expiry notification")
		notificationService := NewNotificationService()
		if err := notificationService.NotifyContractExpiring(); err != nil {
			log.Printf("Error notifying contract expiry: %v", err)
		}
	})

	// Run daily contract auto-termination at 10:00
	go s.scheduleDaily(10*time.Hour, func() {
		log.Println("Running: Auto-terminate expired contracts")
		contractService := NewContractService()
		if err := contractService.AutoTerminateExpiredContracts(); err != nil {
			log.Printf("Error auto-terminating contracts: %v", err)
		}
	})

	// Run daily equipment maintenance check at 11:00
	go s.scheduleDaily(11*time.Hour, func() {
		log.Println("Running: Equipment maintenance notification")
		notificationService := NewNotificationService()
		if err := notificationService.NotifyEquipmentMaintenanceDue(); err != nil {
			log.Printf("Error notifying equipment maintenance: %v", err)
		}
	})

	// Run cleanup of old notifications at 02:00
	go s.scheduleDaily(2*time.Hour, func() {
		log.Println("Running: Cleanup old notifications")
		notificationService := NewNotificationService()
		if err := notificationService.CleanupOldNotifications(); err != nil {
			log.Printf("Error cleaning up notifications: %v", err)
		}
	})
}

// StopScheduler stops all background jobs.
// Dùng non-blocking send để defer khi panic không bị deadlock
// (các goroutine scheduler chỉ select stopChannel sau khi time.Sleep xong).
func (s *SchedulerService) StopScheduler() {
	if !s.isRunning {
		return
	}

	s.isRunning = false
	select {
	case s.stopChannel <- true:
	default:
		// Không ai sẵn sàng nhận, bỏ qua để không block defer.
	}
	log.Println("✓ Scheduler stopped")
}

// scheduleDaily schedules a task to run at a specific time each day
func (s *SchedulerService) scheduleDaily(targetTime time.Duration, task func()) {
	for {
		select {
		case <-s.stopChannel:
			return
		default:
			now := time.Now()
			nextRun := now.Add(targetTime - time.Duration(now.Hour())*time.Hour - time.Duration(now.Minute())*time.Minute - time.Duration(now.Second())*time.Second)

			// If time has already passed today, schedule for tomorrow
			if nextRun.Before(now) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			waitDuration := nextRun.Sub(now)
			time.Sleep(waitDuration)

			if !s.isRunning {
				return
			}

			task()
		}
	}
}

// ScheduleTask schedules a one-time task
func (s *SchedulerService) ScheduleTask(delay time.Duration, task func()) {
	go func() {
		time.Sleep(delay)
		if s.isRunning {
			task()
		}
	}()
}
