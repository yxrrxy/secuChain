document.addEventListener('DOMContentLoaded', () => {
    const body = document.body;
    const particles = [];
    const colors = ['#94E2D5', '#89DCEB', '#BAC2DE', '#B4BEFE'];

    function createParticle(x, y) {
        const particle = document.createElement('div');
        particle.className = 'particle';
        particle.style.left = `${x}px`;
        particle.style.top = `${y}px`;
        particle.style.backgroundColor = colors[Math.floor(Math.random() * colors.length)];
        body.appendChild(particle);
        particles.push(particle);

        setTimeout(() => {
            particle.remove();
            particles.shift();
        }, 1500);
    }

    body.addEventListener('mousemove', (e) => {
        if (Math.random() > 0.5) {  // 降低粒子生成频率
            createParticle(e.clientX, e.clientY);
        }
    });
});