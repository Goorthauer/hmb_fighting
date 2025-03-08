export const canvas = document.getElementById('gameCanvas');
export const ctx = canvas.getContext('2d');
export const cellWidth = 60; // Увеличиваем с 50 до 60
export const cellHeight = 60; // Увеличиваем с 50 до 60

// Функция для установки размеров канваса
function setCanvasSize() {
    const width = window.innerWidth;
    if (width <= 768) {
        canvas.width = 16 * cellWidth * 0.75;  // 720px (75% от 960px)
        canvas.height = 9 * cellHeight * 0.75; // 405px (75% от 540px)
    } else if (width <= 1200) {
        canvas.width = 16 * cellWidth * 0.875;  // 840px (87.5% от 960px)
        canvas.height = 9 * cellHeight * 0.875; // 472.5px (87.5% от 540px)
    } else {
        canvas.width = 16 * cellWidth;  // 960px для 16 клеток
        canvas.height = 9 * cellHeight; // 540px для 9 клеток
    }
}

// Устанавливаем начальные размеры
setCanvasSize();

// Обновляем размеры при изменении окна
window.addEventListener('resize', setCanvasSize);