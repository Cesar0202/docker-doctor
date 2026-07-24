document.addEventListener('DOMContentLoaded', () => {
    // Animaciones al hacer scroll (Intersection Observer)
    const fadeElements = document.querySelectorAll('.fade-in, .feature-card, .install-card');
    
    // Add fade-in class to cards if they don't have it
    document.querySelectorAll('.feature-card, .install-card').forEach(el => {
        if (!el.classList.contains('fade-in')) {
            el.classList.add('fade-in');
        }
    });

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
                // Optional: stop observing once it's visible
                // observer.unobserve(entry.target);
            }
        });
    }, {
        threshold: 0.1,
        rootMargin: "0px 0px -50px 0px"
    });

    // Trigger visible immediately for hero elements without intersection observer
    // to avoid flickering on initial load
    setTimeout(() => {
        document.querySelectorAll('.hero .fade-in').forEach(el => {
            el.classList.add('visible');
        });
    }, 100);

    fadeElements.forEach(element => {
        if (!element.closest('.hero')) {
            observer.observe(element);
        }
    });

    // Lógica para copiar comandos al portapapeles
    const copyButtons = document.querySelectorAll('.copy-btn');
    
    copyButtons.forEach(btn => {
        btn.addEventListener('click', (e) => {
            const targetId = e.target.getAttribute('data-target');
            const textToCopy = document.getElementById(targetId).textContent;
            
            navigator.clipboard.writeText(textToCopy).then(() => {
                const originalText = e.target.textContent;
                e.target.textContent = '¡Copiado!';
                e.target.style.background = 'rgba(0, 255, 136, 0.2)';
                e.target.style.color = '#00ff88';
                
                setTimeout(() => {
                    e.target.textContent = originalText;
                    e.target.style.background = '';
                    e.target.style.color = '';
                }, 2000);
            }).catch(err => {
                console.error('Error al copiar: ', err);
            });
        });
    });
});
