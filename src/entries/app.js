import '../stylesheets/styles.css'; 
import '../stylesheets/dashboard.css'; 
import { AutoComplete } from '../plugins/AutoComplete.js';
import '../plugins/AutoComplete.css';
import { FileUpload } from '../plugins/FileUpload.js';
import '../plugins/FileUpload.css';

window.AutoComplete = AutoComplete;
window.FileUpload = FileUpload;

document.addEventListener('DOMContentLoaded', function() {
  const sidebarToggle = document.getElementById('sidebarToggle');
  const sidebar = document.querySelector('.sidebar');
  const mainContent = document.querySelector('.main-content');
  
  // Verificar el estado inicial en localStorage
  let isCollapsed = localStorage.getItem('sidebarCollapsed') === 'true';
  
  // Función para actualizar el estado
  function updateSidebarState() {
    if (window.innerWidth <= 992) {
      // En móvil: siempre empezar colapsado
      sidebar.classList.add('collapsed');
    } else {
      // En desktop: usar el estado guardado
      sidebar.classList.toggle('collapsed', isCollapsed);
    }
  }
  
  // Aplicar estado inicial
  updateSidebarState();
  
  // Evento click del toggle
  sidebarToggle.addEventListener('click', function() {
    if (window.innerWidth <= 992) {
      // En móvil: toggle simple
      sidebar.classList.toggle('collapsed');
    } else {
      // En desktop: guardar preferencia
      isCollapsed = !isCollapsed;
      localStorage.setItem('sidebarCollapsed', isCollapsed);
      sidebar.classList.toggle('collapsed');
    }
  });
  
  // Cerrar sidebar al hacer clic en enlaces (solo móvil)
  document.querySelectorAll('.sidebar .nav-link').forEach(link => {
    link.addEventListener('click', function() {
      if (window.innerWidth <= 992) {
        sidebar.classList.add('collapsed');
      }
    });
  });
  
  // Actualizar al cambiar tamaño de pantalla
  window.addEventListener('resize', function() {
    updateSidebarState();
  });
});