export const canvas = document.getElementById('gameCanvas');
export const ctx = canvas.getContext('2d');
export const cellWidth = 50;
export const cellHeight = 50;

// Функция для установки размеров канваса
function setCanvasSize() {
    const width = window.innerWidth;
    if (width <= 768) {
        canvas.width = 600; // Соответствует @media (max-width: 768px)
        canvas.height = 300;
    } else if (width <= 1200) {
        canvas.width = 800; // Соответствует @media (max-width: 1200px)
        canvas.height = 400;
    } else {
        canvas.width = 1000; // Базовый размер
        canvas.height = 500;
    }
}

// Устанавливаем начальные размеры
setCanvasSize();

// Обновляем размеры при изменении окна
window.addEventListener('resize', setCanvasSize);