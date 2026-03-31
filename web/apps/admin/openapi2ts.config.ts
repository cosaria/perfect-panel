const config = [
  {
    requestLibPath: "import request from '@/utils/request';",
    schemaPath:
      'https://raw.githubusercontent.com/perfect-panel/ppanel-docs/refs/heads/main/public/swagger/common.json',
    serversPath: './services',
    projectName: 'common',
  },
  {
    requestLibPath: "import request from '@/utils/request';",
    schemaPath:
      'https://raw.githubusercontent.com/perfect-panel/ppanel-docs/refs/heads/main/public/swagger/admin.json',
    serversPath: './services',
    projectName: 'admin',
  },
];

export default config;
