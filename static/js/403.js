document.addEventListener('DOMContentLoaded', function() {
  const backBtn = document.getElementById('backBtn');
  
  backBtn.addEventListener('click', function(e) {
    e.preventDefault();
    
    // Verificar si hay historial previo
    if (window.history.length > 1) {
      // Regresar a la página anterior
      window.history.back();
    } else {
      // Si no hay historial, redirigir al inicio
      window.location.href = URLS.BASE + "/";
    }
  });
});