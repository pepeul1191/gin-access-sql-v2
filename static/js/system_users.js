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

    // Crear objeto con todos los datos a enviar
    const formData = {
      employees: employeesData,
    };

    // Enviar datos al servidor
    fetch(`/systems/${SYSTEM_ID}/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRFToken': '{{ csrf_token }}'
      },
      body: JSON.stringify(formData)
    })
    .then(response => {
      if (response.ok) {
        //window.location.reload(); // Recargar la página después de guardar
        alert('Cambios guardados correctamente');
      } else {
        alert('Error al guardar los cambios');
      }
    })
    .catch(error => {
      console.error('Error:', error);
      alert('Error al conectar con el servidor');
    });
  });
});