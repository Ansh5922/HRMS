-- Insert system permissions
INSERT INTO permissions (resource, action) VALUES
  ('employee','read'),   ('employee','write'),   ('employee','delete'),
  ('attendance','read'), ('attendance','write'),
  ('leave','read'),      ('leave','write'),       ('leave','approve'),
  ('payroll','read'),    ('payroll','write'),      ('payroll','approve'),
  ('recruitment','read'),('recruitment','write'),
  ('performance','read'),('performance','write'),
  ('training','read'),   ('training','write'),
  ('reports','read'),
  ('settings','read'),   ('settings','write'),
  ('users','read'),      ('users','write'),        ('users','delete')
ON CONFLICT (resource, action) DO NOTHING;
