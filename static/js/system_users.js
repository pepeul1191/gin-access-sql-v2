document.addEventListener('DOMContentLoaded', function() {
  // Checkbox para seleccionar/deseleccionar todos
  const selectAllCheckbox = document.getElementById('select-all');
  const employeeCheckboxes = document.querySelectorAll('.user-checkbox');

  // Evento para el checkbox "Seleccionar todos"
  selectAllCheckbox.addEventListener('change', function() {
    employeeCheckboxes.forEach(checkbox => {
      checkbox.checked = selectAllCheckbox.checked;
    });
  });
  
  // Evento para los checkboxes individuales
  employeeCheckboxes.forEach(checkbox => {
    checkbox.addEventListener('change', function() {
      if (!this.checked) {
        selectAllCheckbox.checked = false;
      } else {
        // Verificar si todos los checkboxes están seleccionados
        const allChecked = Array.from(employeeCheckboxes).every(cb => cb.checked);
        selectAllCheckbox.checked = allChecked;
      }
    });
  });

  // Botón Guardar Cambios
  document.querySelector('.btn-success').addEventListener('click', function(e) {
    e.preventDefault();
    // Crear lista de employees con su estado de selección
    const employeesData = Array.from(employeeCheckboxes).map(checkbox => {
      return {
        id: parseInt(checkbox.value),
        selected: checkbox.checked
      };
    });


    // Enviar datos al servidor
    fetch(`/systems/${SYSTEM_ID}/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRFToken': CSRF_TOKEN
      },
      body: JSON.stringify(employeesData)
    })
    .then(response => {
      if (response.ok) {
        //window.location.reload(); // Recargar la página después de guardar
        window.location.href = `/systems/${SYSTEM_ID}/users?message=Cambios%20guardados%20correctamente&type=success`
      } else {
        window.location.href = `/systems/${SYSTEM_ID}/users?message=Error%20al%20guardar%20los%20cambios&type=danger`
      }
    })
    .catch(error => {
      window.location.href = `/systems/${SYSTEM_ID}/users?message=Error%20al%20conectar%20con%20el%20servidor&type=danger`
      console.error('Error:', error);
    });
  });
});