-- migrate:up

INSERT INTO systems (id, name, description, repository, created, updated) VALUES 
(3, 'Inventory Management System', 'Sistema para gestionar inventario y stock de productos.', 'https://github.com/company/inventory-system', '2024-01-15 09:30:00', '2024-06-10 14:20:00'),
(4, 'Customer Relationship Manager', 'Plataforma integral para gestión de clientes y ventas.', 'https://github.com/company/crm-system', '2023-11-20 11:45:00', '2024-05-22 10:15:00'),
(5, 'Financial Analytics Suite', 'Herramientas avanzadas para análisis financiero y reporting.', 'https://github.com/company/finance-analytics', '2024-02-08 08:20:00', '2024-06-05 16:40:00'),
(6, 'HR Management Platform', 'Sistema completo para gestión de recursos humanos.', 'https://github.com/company/hr-platform', '2023-09-12 13:10:00', '2024-04-18 09:25:00'),
(7, 'E-commerce Solution', 'Plataforma de comercio electrónico con múltiples funcionalidades.', 'https://github.com/company/ecommerce-solution', '2024-03-03 10:30:00', '2024-06-12 11:50:00'),
(8, 'Project Tracking Tool', 'Sistema para seguimiento y gestión de proyectos.', 'https://github.com/company/project-tracker', '2023-12-05 14:15:00', '2024-05-28 15:30:00'),
(9, 'Data Visualization Dashboard', 'Panel interactivo para visualización de datos empresariales.', 'https://github.com/company/data-dashboard', '2024-01-28 16:40:00', '2024-06-08 12:20:00'),
(10, 'Document Management System', 'Sistema centralizado para gestión documental.', 'https://github.com/company/document-manager', '2023-10-17 09:50:00', '2024-04-25 14:35:00'),
(11, 'Supply Chain Optimizer', 'Herramienta para optimización de cadena de suministro.', 'https://github.com/company/supply-chain', '2024-02-22 11:25:00', '2024-05-15 10:45:00'),
(12, 'Business Intelligence Suite', 'Suite completa de inteligencia de negocios.', 'https://github.com/company/bi-suite', '2023-08-30 15:20:00', '2024-03-20 13:10:00'),
(13, 'Marketing Automation Platform', 'Plataforma para automatización de campañas de marketing.', 'https://github.com/company/marketing-automation', '2024-03-17 12:30:00', '2024-06-14 09:40:00'),
(14, 'Quality Control System', 'Sistema para control de calidad y garantía.', 'https://github.com/company/quality-control', '2023-11-05 10:15:00', '2024-04-12 16:25:00'),
(15, 'Asset Management Tool', 'Herramienta para gestión de activos empresariales.', 'https://github.com/company/asset-manager', '2024-01-08 14:50:00', '2024-05-05 11:30:00'),
(16, 'Time Tracking Application', 'Aplicación para registro y control de tiempo.', 'https://github.com/company/time-tracker', '2023-10-03 08:45:00', '2024-03-28 14:15:00'),
(17, 'Compliance Management System', 'Sistema para gestión de cumplimiento normativo.', 'https://github.com/company/compliance-system', '2024-02-14 13:20:00', '2024-06-03 10:50:00'),
(18, 'Customer Support Portal', 'Portal integral para atención al cliente.', 'https://github.com/company/support-portal', '2023-09-25 11:10:00', '2024-04-08 15:40:00'),
(19, 'Business Process Manager', 'Herramienta para modelado de procesos de negocio.', 'https://github.com/company/bpm-tool', '2024-03-08 16:30:00', '2024-06-11 12:25:00'),
(20, 'Risk Assessment Platform', 'Plataforma para evaluación de riesgos empresariales.', 'https://github.com/company/risk-assessment', '2023-12-18 09:35:00', '2024-05-20 14:50:00'),
(21, 'Performance Monitoring System', 'Sistema para monitoreo de rendimiento de aplicaciones.', 'https://github.com/company/performance-monitor', '2024-01-22 10:40:00', '2024-05-12 11:15:00'),
(22, 'Collaboration Workspace', 'Espacio de trabajo colaborativo para equipos.', 'https://github.com/company/collaboration-space', '2023-08-15 14:25:00', '2024-02-28 16:30:00'),
(23, 'Data Integration Hub', 'Centralizador para integración de datos heterogéneos.', 'https://github.com/company/data-integration', '2024-02-05 12:15:00', '2024-05-25 09:45:00'),
(24, 'Mobile Workforce Management', 'Sistema para gestión de fuerza laboral móvil.', 'https://github.com/company/mobile-workforce', '2023-11-28 15:50:00', '2024-04-22 13:20:00'),
(25, 'Expense Reporting System', 'Sistema automatizado para reporte de gastos.', 'https://github.com/company/expense-reporter', '2024-03-25 08:30:00', '2024-06-13 14:10:00'),
(26, 'Learning Management System', 'Plataforma para gestión de capacitación y desarrollo.', 'https://github.com/company/lms-platform', '2023-10-10 13:45:00', '2024-03-15 10:25:00'),
(27, 'Business Continuity Planner', 'Herramienta para planificación de continuidad operativa.', 'https://github.com/company/business-continuity', '2024-01-03 16:20:00', '2024-05-18 15:35:00'),
(28, 'Contract Management System', 'Sistema para gestión integral de contratos.', 'https://github.com/company/contract-manager', '2023-12-12 11:05:00', '2024-04-30 12:40:00'),
(29, 'Data Privacy Compliance', 'Herramienta para cumplimiento de privacidad de datos.', 'https://github.com/company/data-privacy', '2024-02-28 14:55:00', '2024-06-07 08:50:00'),
(30, 'Enterprise Search Engine', 'Motor de búsqueda corporativo unificado.', 'https://github.com/company/enterprise-search', '2023-09-18 10:20:00', '2024-03-10 16:45:00');

-- migrate:down

DELETE FROM systems WHERE id >= 3;