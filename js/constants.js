export const canvas = document.getElementById('gameCanvas');
export const ctx = canvas.getContext('2d');
export const cellWidth = 50;
export const cellHeight = 50;

canvas.width = 20 * cellWidth; // Устанавливаем размер canvas
canvas.height = 10 * cellHeight;