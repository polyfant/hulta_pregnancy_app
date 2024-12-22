document.addEventListener('DOMContentLoaded', function() {
    // Update progress circle
    const progressCircle = document.querySelector('.progress-circle');
    const progressValue = progressCircle.querySelector('.progress-value').textContent;
    const progress = parseFloat(progressValue);
    progressCircle.style.setProperty('--progress', `${progress * 3.6}deg`);

    // Add touch feedback for mobile
    const scheduleItems = document.querySelectorAll('.schedule-item');
    scheduleItems.forEach(item => {
        item.addEventListener('touchstart', function() {
            this.style.transform = 'scale(0.95)';
        });
        
        item.addEventListener('touchend', function() {
            this.style.transform = 'scale(1)';
        });
    });

    // Auto-refresh data every minute
    setInterval(async function() {
        try {
            const response = await fetch('/api/pregnancy/status');
            if (response.ok) {
                const data = await response.json();
                updatePregnancyData(data);
            }
        } catch (error) {
            console.error('Failed to update pregnancy data:', error);
        }
    }, 60000);
});

function updatePregnancyData(data) {
    // Update progress
    const progressValue = document.querySelector('.progress-value');
    const progressCircle = document.querySelector('.progress-circle');
    progressValue.textContent = `${data.Progress.toFixed(2)}%`;
    progressCircle.style.setProperty('--progress', `${data.Progress * 3.6}deg`);

    // Update details
    document.querySelector('.detail-item:nth-child(1) .value').textContent = data.CurrentDay;
    document.querySelector('.detail-item:nth-child(2) .value').textContent = data.RemainingDays;
    
    const stageElement = document.querySelector('.detail-item:nth-child(3) .value');
    stageElement.textContent = data.Stage;
    stageElement.className = `value stage-${data.Stage}`;

    // Update monitoring schedule
    const scheduleItems = document.querySelectorAll('.schedule-item');
    scheduleItems[0].classList.toggle('active', data.Schedule.TemperatureCheck);
    scheduleItems[1].classList.toggle('active', data.Schedule.BehaviorCheck);
    scheduleItems[2].classList.toggle('active', data.Schedule.UdderCheck);
    scheduleItems[3].classList.toggle('active', data.Schedule.VulvaCheck);

    // Update check frequency
    const checkFrequency = document.querySelector('.check-frequency span');
    checkFrequency.textContent = `Check every ${data.Schedule.CheckFrequency} hours`;
    checkFrequency.className = `priority-${data.Schedule.Priority}`;
}
