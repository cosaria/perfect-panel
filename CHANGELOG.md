<a name="readme-top"></a>
# Changelog

## [1.6.3](https://github.com/perfect-panel/ppanel-web/compare/v1.6.2...v1.6.3) (2025-12-08)


### 🐛 Bug Fixes

* **docker**: Update Dockerfiles to create non-root user with proper permissions ([1bfebb6](https://github.com/perfect-panel/ppanel-web/commit/1bfebb6))

## [1.6.2](https://github.com/perfect-panel/ppanel-web/compare/v1.6.1...v1.6.2) (2025-12-08)


### 🐛 Bug Fixes

* **package**: Update dependencies and upgrade React and Next.js versions. ([7d0866e](https://github.com/perfect-panel/ppanel-web/commit/7d0866e))

## [1.6.1](https://github.com/perfect-panel/ppanel-web/compare/v1.6.0...v1.6.1) (2025-11-05)


### 🐛 Bug Fixes

* Fixing issues with generating standard and quantum-resistant encryption keys ([5eac6a9](https://github.com/perfect-panel/ppanel-web/commit/5eac6a9))

<a name="readme-top"></a>

# Changelog

# [1.6.0](https://github.com/perfect-panel/ppanel-web/compare/v1.5.4...v1.6.0) (2025-10-28)

### ✨ Features

- Add server installation dialog and commands ([4429c9d](https://github.com/perfect-panel/ppanel-web/commit/4429c9d))

### 🐛 Bug Fixes

- Add typeRoots configuration to ensure type definitions are resolved correctly ([ad60ea9](https://github.com/perfect-panel/ppanel-web/commit/ad60ea9))

<a name="readme-top"></a>

# Changelog

## [1.5.4](https://github.com/perfect-panel/ppanel-web/compare/v1.5.3...v1.5.4) (2025-10-26)

### 🐛 Bug Fixes

- Update generateRealityKeyPair to use async key generation ([e60e369](https://github.com/perfect-panel/ppanel-web/commit/e60e369))
- Update the wallet localization file and add new fields such as automatic reset and recharge ([88aa965](https://github.com/perfect-panel/ppanel-web/commit/88aa965))

## [1.5.3](https://github.com/perfect-panel/ppanel-web/compare/v1.5.2...v1.5.3) (2025-10-21)

### 🐛 Bug Fixes

- Fix bugs ([a46657d](https://github.com/perfect-panel/ppanel-web/commit/a46657d))
- Fix dependencies ([8bd25d6](https://github.com/perfect-panel/ppanel-web/commit/8bd25d6))
- Remove unnecessary migration function code and add device configuration options ([521a7a9](https://github.com/perfect-panel/ppanel-web/commit/521a7a9))
- Update bun.lockb to reflect dependency changes ([ca892dd](https://github.com/perfect-panel/ppanel-web/commit/ca892dd))

<a name="readme-top"></a>

# Changelog

## [1.5.2](https://github.com/perfect-panel/ppanel-web/compare/v1.5.1...v1.5.2) (2025-09-29)

### 🐛 Bug Fixes

- Add step attribute to datetime-local inputs for precise time selection in forms ([32fd181](https://github.com/perfect-panel/ppanel-web/commit/32fd181))
- Rename 'hysteria2' to 'hysteria' across protocol definitions and schemas for consistency ([5816dd5](https://github.com/perfect-panel/ppanel-web/commit/5816dd5))
- Update protocol options in ServerConfig for accuracy and consistency ([9266529](https://github.com/perfect-panel/ppanel-web/commit/9266529))

## [1.5.1](https://github.com/perfect-panel/ppanel-web/compare/v1.5.0...v1.5.1) (2025-09-28)

### 🐛 Bug Fixes

- Simplify protocol enable checks by removing unnecessary false comparisons ([4828700](https://github.com/perfect-panel/ppanel-web/commit/4828700))

# [1.5.0](https://github.com/perfect-panel/ppanel-web/compare/v1.4.8...v1.5.0) (2025-09-28)

### ✨ Features

- Update server configuration translations for multiple languages ([fc43de1](https://github.com/perfect-panel/ppanel-web/commit/fc43de1))

### 🐛 Bug Fixes

- Add DynamicMultiplier component for managing node multipliers and update ServersPage layout ([bb6671c](https://github.com/perfect-panel/ppanel-web/commit/bb6671c))
- Remove unnecessary blank lines in multiple index files for cleaner code structure ([6a823b8](https://github.com/perfect-panel/ppanel-web/commit/6a823b8))
- Remove unused ratio variable from server traffic log and server form for cleaner code ([55034dc](https://github.com/perfect-panel/ppanel-web/commit/55034dc))
- Update Badge variants and restructure traffic ratio display in ServersPage ([3d778e5](https://github.com/perfect-panel/ppanel-web/commit/3d778e5))
- Update minimum ratio value to 0 in protocol fields and adjust related schemas; enhance unit conversion in ServerConfig ([3b6ef17](https://github.com/perfect-panel/ppanel-web/commit/3b6ef17))
- Update protocol fields to use 'obfs' instead of 'security' and adjust related configurations ([4abdd36](https://github.com/perfect-panel/ppanel-web/commit/4abdd36))

<a name="readme-top"></a>

# Changelog

## [1.4.8](https://github.com/perfect-panel/ppanel-web/compare/v1.4.7...v1.4.8) (2025-09-23)

### 🐛 Bug Fixes

- Rename 'server_id' to 'protocol' in NodesPage and clean up unused imports and code in ServerConfig ([70b3484](https://github.com/perfect-panel/ppanel-web/commit/70b3484))
- Update announcement page to display timeline of announcements with Markdown content ([3c036eb](https://github.com/perfect-panel/ppanel-web/commit/3c036eb))
- Update Empty component to support border prop and adjust usage in various pages ([ce9ab89](https://github.com/perfect-panel/ppanel-web/commit/ce9ab89))

## [1.4.7](https://github.com/perfect-panel/ppanel-web/compare/v1.4.6...v1.4.7) (2025-09-23)

### 🐛 Bug Fixes

- Add unique key to ProTable for improved rendering with user ID filters ([2bff15f](https://github.com/perfect-panel/ppanel-web/commit/2bff15f))
- Adjust layout spacing and chart aspect ratio in ServerConfig component ([05a61d8](https://github.com/perfect-panel/ppanel-web/commit/05a61d8))
- Refactor server ID cell rendering for improved readability and consistency ([0345b7c](https://github.com/perfect-panel/ppanel-web/commit/0345b7c))
- Update announcement page to format creation date and enhance content display ([8445e30](https://github.com/perfect-panel/ppanel-web/commit/8445e30))
- Update OnlineUsersCell to display user count with icon instead of badge ([7a4ebdf](https://github.com/perfect-panel/ppanel-web/commit/7a4ebdf))
- Update subscribe name fallback to return '--' instead of 'Unknown' ([0a07d25](https://github.com/perfect-panel/ppanel-web/commit/0a07d25))

## [1.4.6](https://github.com/perfect-panel/ppanel-web/compare/v1.4.5...v1.4.6) (2025-09-17)

### 🎫 Chores

- Merge branch 'main' into develop ([41f06bf](https://github.com/perfect-panel/ppanel-web/commit/41f06bf))

### 🐛 Bug Fixes

- Add loaded state to node, server, and subscribe stores for better loading management ([13dce0c](https://github.com/perfect-panel/ppanel-web/commit/13dce0c))
- Removed node metadata fields in subscription schema to simplify structure ([0cadd83](https://github.com/perfect-panel/ppanel-web/commit/0cadd83))

## [1.4.5](https://github.com/perfect-panel/ppanel-web/compare/v1.4.4...v1.4.5) (2025-09-17)

### ♻ Code Refactoring

- Replace useQuery with Zustand store for subscription and node data management ([c6dd0b6](https://github.com/perfect-panel/ppanel-web/commit/c6dd0b6))
- Simplify TemplatePreview component structure by consolidating Sheet and Button elements ([1b715c5](https://github.com/perfect-panel/ppanel-web/commit/1b715c5))

### 🐛 Bug Fixes

- Add showLineNumbers prop handling in MonacoEditor for improved placeholder positioning ([bd67ece](https://github.com/perfect-panel/ppanel-web/commit/bd67ece))
- Add fetchTags method to NodesPage and ensure tags are fetched alongside nodes ([a3c5e31](https://github.com/perfect-panel/ppanel-web/commit/a3c5e31))
- Add NEXT_PUBLIC_HIDDEN_TUTORIAL_DOCUMENT to control tutorial visibility and update page query accordingly ([e94405d](https://github.com/perfect-panel/ppanel-web/commit/e94405d))
- Add subscribeSchema for subscription management with detailed proxy and user information ([49b3dcc](https://github.com/perfect-panel/ppanel-web/commit/49b3dcc))
- Enhance server ID display in ServerTrafficLogPage with badges and server ratio ([6dfac27](https://github.com/perfect-panel/ppanel-web/commit/6dfac27))
- Update platform handling in Content component to ensure available platforms are correctly filtered and displayed ([1dde708](https://github.com/perfect-panel/ppanel-web/commit/1dde708))

<a name="readme-top"></a>

# Changelog

## [1.4.4](https://github.com/perfect-panel/ppanel-web/compare/v1.4.3...v1.4.4) (2025-09-16)

### 🐛 Bug Fixes

- Add minimum value constraint for count and user limit inputs in CouponForm ([6991b69](https://github.com/perfect-panel/ppanel-web/commit/6991b69))
- Add protocol-related constants, default configurations, field definitions, and validation modes ([d685407](https://github.com/perfect-panel/ppanel-web/commit/d685407))
- Added the enabled field in the protocol configuration, updated the related type definition and default configuration ([2b0cf9a](https://github.com/perfect-panel/ppanel-web/commit/2b0cf9a))
- Filter available protocols to exclude disabled ones in NodeForm ([982d288](https://github.com/perfect-panel/ppanel-web/commit/982d288))
- Refactor key generation logic and update dependencies for ML-KEM-768 integration ([b8f630f](https://github.com/perfect-panel/ppanel-web/commit/b8f630f))

<a name="readme-top"></a>

# Changelog

## [1.4.3](https://github.com/perfect-panel/ppanel-web/compare/v1.4.2...v1.4.3) (2025-09-16)

### 🐛 Bug Fixes

- Add success toast message for sorting in nodes and servers pages ([2d5175d](https://github.com/perfect-panel/ppanel-web/commit/2d5175d))
- Implement encryption and obfuscation features in protocol configuration ([54de16b](https://github.com/perfect-panel/ppanel-web/commit/54de16b))
- Refactor toB64 function to toB64Url for URL-safe base64 encoding in VlessX25519Pair generation ([8700cf6](https://github.com/perfect-panel/ppanel-web/commit/8700cf6))
- Simplify initialValues assignment and update node submission logic in NodesPage ([05d6c89](https://github.com/perfect-panel/ppanel-web/commit/05d6c89))
- Update bun.lockb to reflect dependency changes ([ebcebd7](https://github.com/perfect-panel/ppanel-web/commit/ebcebd7))

<a name="readme-top"></a>

# Changelog

## [1.4.2](https://github.com/perfect-panel/ppanel-web/compare/v1.4.1...v1.4.2) (2025-09-15)

### 🐛 Bug Fixes

- Add GitHub template repository link to ProtocolForm header for easier access ([8a0baf3](https://github.com/perfect-panel/ppanel-web/commit/8a0baf3))
- Add readOnly prop to MonacoEditor and TemplatePreview components for improved content handling ([c4c4d5a](https://github.com/perfect-panel/ppanel-web/commit/c4c4d5a))

## [1.4.1](https://github.com/perfect-panel/ppanel-web/compare/v1.4.0...v1.4.1) (2025-09-15)

### 🐛 Bug Fixes

- Add copy subscription functionality to user subscription dashboard and improve localization for new features ([e2a357f](https://github.com/perfect-panel/ppanel-web/commit/e2a357f))
- Simplify handleInputBlur function by removing unnecessary setTimeout ([6341562](https://github.com/perfect-panel/ppanel-web/commit/6341562))
- Update TemplatePreview to use MonacoEditor for content display and improve error handling ([1d52642](https://github.com/perfect-panel/ppanel-web/commit/1d52642))
- Update user dashboard link to correct path ([131693b](https://github.com/perfect-panel/ppanel-web/commit/131693b))

<a name="readme-top"></a>

# Changelog

# [1.4.0](https://github.com/perfect-panel/ppanel-web/compare/v1.3.0...v1.4.0) (2025-09-14)

### ♻ Code Refactoring

- **logs**: Add localization files and update existing translations for multiple languages ([2f20ac9](https://github.com/perfect-panel/ppanel-web/commit/2f20ac9))
- **subscribe-form**: Replace server_group and server with node_tags and nodes in default values and form schema ([38dda84](https://github.com/perfect-panel/ppanel-web/commit/38dda84))
- Add localization updates for log and server files across multiple languages ([2bcd4cf](https://github.com/perfect-panel/ppanel-web/commit/2bcd4cf))
- Clean up NodeForm and ServerForm components by removing unused functions and optimizing state management ([10250d9](https://github.com/perfect-panel/ppanel-web/commit/10250d9))
- Enhance log pages with Badge component and update translations ([a4de9df](https://github.com/perfect-panel/ppanel-web/commit/a4de9df))
- Make 'scheme' field optional in client form schema ([f4a1237](https://github.com/perfect-panel/ppanel-web/commit/f4a1237))
- Refactor server management API endpoints and typings ([4f7cc80](https://github.com/perfect-panel/ppanel-web/commit/4f7cc80))
- Remove application management forms and related configurations ([0c43844](https://github.com/perfect-panel/ppanel-web/commit/0c43844))
- Remove default options from TagInput component for improved flexibility ([6a3bb70](https://github.com/perfect-panel/ppanel-web/commit/6a3bb70))
- Remove unused preview state variables and add sort order to node properties ([e63f823](https://github.com/perfect-panel/ppanel-web/commit/e63f823))
- Rename buildScheme to buildSchema and update imports in server form components ([ee98e7e](https://github.com/perfect-panel/ppanel-web/commit/ee98e7e))
- Simplify node display in subscribe form and remove unused Badge import ([551135d](https://github.com/perfect-panel/ppanel-web/commit/551135d))
- Update bun.lockb to reflect dependency changes ([ba2b50e](https://github.com/perfect-panel/ppanel-web/commit/ba2b50e))
- Update component imports and improve code consistency ([59faeab](https://github.com/perfect-panel/ppanel-web/commit/59faeab))
- Update dependencies and improve code consistency across multiple files ([e37ae49](https://github.com/perfect-panel/ppanel-web/commit/e37ae49))
- Update localization files and service imports ([d4b37e4](https://github.com/perfect-panel/ppanel-web/commit/d4b37e4))

### ✨ Features

- **config**: Add translations for server configuration in multiple languages ([f9a7ece](https://github.com/perfect-panel/ppanel-web/commit/f9a7ece))
- **logs**: Add various log pages for tracking user activities and system events ([d85af49](https://github.com/perfect-panel/ppanel-web/commit/d85af49))
- Add bandwidth fields and placeholders for upload and download in server configuration forms; update localization files for multiple languages ([3e5402f](https://github.com/perfect-panel/ppanel-web/commit/3e5402f))
- Add batch delete functionality and enhance chart tooltips in statistics cards ([c4f536e](https://github.com/perfect-panel/ppanel-web/commit/c4f536e))
- Add language support and descriptions in product localization files ([fd48856](https://github.com/perfect-panel/ppanel-web/commit/fd48856))
- Add log cleanup settings and update localization files ([6ccf9b8](https://github.com/perfect-panel/ppanel-web/commit/6ccf9b8))
- Add queryNodeTag function and integrate tag retrieval in NodeForm and SubscribeForm components ([4563c57](https://github.com/perfect-panel/ppanel-web/commit/4563c57))
- Add quota management features and localization updates ([fce627b](https://github.com/perfect-panel/ppanel-web/commit/fce627b))
- Add server form component with protocol configuration and localization support ([217ddce](https://github.com/perfect-panel/ppanel-web/commit/217ddce))
- Enhance TagInput component with option handling and improved tag addition logic ([b6e778d](https://github.com/perfect-panel/ppanel-web/commit/b6e778d))
- Implement data migration functionality and update localization files ([6d81bfd](https://github.com/perfect-panel/ppanel-web/commit/6d81bfd))
- Refactor user detail and subscription management components ([973c06f](https://github.com/perfect-panel/ppanel-web/commit/973c06f))
- Update server list fetching logic and adjust query parameters ([5272360](https://github.com/perfect-panel/ppanel-web/commit/5272360))

### 🐛 Bug Fixes

- **page**: Refine version checking logic and remove unnecessary comments for clarity ([26176a7](https://github.com/perfect-panel/ppanel-web/commit/26176a7))
- Add localization updates and new utility functions ([4da5960](https://github.com/perfect-panel/ppanel-web/commit/4da5960))
- Add user_subscribe_id filter to SubscribeLogPage and update typings; disable eslint in service index files ([ab6f6a6](https://github.com/perfect-panel/ppanel-web/commit/ab6f6a6))
- Added system version card and system log dialog components; updated statistics page to include total server and user statistics ([fe69980](https://github.com/perfect-panel/ppanel-web/commit/fe69980))
- Adjust timestamp calculations to use milliseconds instead of seconds in quota broadcast form submission ([8dffd69](https://github.com/perfect-panel/ppanel-web/commit/8dffd69))
- Correct cookie key format for sidebar state retrieval ([9e01f4f](https://github.com/perfect-panel/ppanel-web/commit/9e01f4f))
- Improve UI for protocol status display in server form component ([461fdb1](https://github.com/perfect-panel/ppanel-web/commit/461fdb1))
- Increase pagination size limit for server and node lists to improve data retrieval ([7c0a312](https://github.com/perfect-panel/ppanel-web/commit/7c0a312))
- Optimize Task Manager to use milliseconds to calculate timers, and simplify import statements ([d8fa13b](https://github.com/perfect-panel/ppanel-web/commit/d8fa13b))
- Refactor protocol status display for improved readability in server form component ([f700be0](https://github.com/perfect-panel/ppanel-web/commit/f700be0))
- Remove GroupTable and related components, simplify SubscribeTable and update language handling in subscription forms ([1ab9b39](https://github.com/perfect-panel/ppanel-web/commit/1ab9b39))
- Remove redundant transport label from localization files ([88ce8d7](https://github.com/perfect-panel/ppanel-web/commit/88ce8d7))
- Remove unnecessary comments and improve variable handling in GoTemplateEditor; ensure zero value displays as empty in EnhancedInput ([c4a47a4](https://github.com/perfect-panel/ppanel-web/commit/c4a47a4))
- Remove unnecessary comments to simplify code readability ([a988cb3](https://github.com/perfect-panel/ppanel-web/commit/a988cb3))
- Replace MarkdownEditor with HTMLEditor in EmailBroadcastForm; simplify content rendering in EmailTaskManager; disable eslint in service index files ([e2d83ec](https://github.com/perfect-panel/ppanel-web/commit/e2d83ec))
- Replace redundant icon rendering with a single instance in system version card component ([e4429a5](https://github.com/perfect-panel/ppanel-web/commit/e4429a5))
- Update billing URL fetching logic and improve version handling in system version card component ([1c8d4af](https://github.com/perfect-panel/ppanel-web/commit/1c8d4af))
- Update condition for plugin field in PROTOCOL_FIELDS to include specific plugins ([6376ec1](https://github.com/perfect-panel/ppanel-web/commit/6376ec1))
- Update getAppSubLink function to improve URL handling and encoding logic ([351fffc](https://github.com/perfect-panel/ppanel-web/commit/351fffc))
- Update order_id to order_no in BalanceLogPage and related typings; enhance timezone switch component with additional features and localization updates ([ac36075](https://github.com/perfect-panel/ppanel-web/commit/ac36075))
- Update protocol plugin handling and add new options in typings ([6ca2433](https://github.com/perfect-panel/ppanel-web/commit/6ca2433))
- Update release URLs to include 'v' prefix for versioning in system version card component ([1d526b5](https://github.com/perfect-panel/ppanel-web/commit/1d526b5))
- Update SidebarLeft component styles to enhance hover effects ([e4fbd5c](https://github.com/perfect-panel/ppanel-web/commit/e4fbd5c))
- Update validation for days and gift_value fields in QuotaBroadcastForm; set default values to avoid errors ([39d746f](https://github.com/perfect-panel/ppanel-web/commit/39d746f))

### 📝 Documentation

- Merge branch 'main' into develop ([d7f8b3b](https://github.com/perfect-panel/ppanel-web/commit/d7f8b3b))

<a name="readme-top"></a>

# Changelog

# [1.3.0](https://github.com/perfect-panel/ppanel-web/compare/v1.2.0...v1.3.0) (2025-08-15)

### ♻ Code Refactoring

- Refactoring and adding multiple features ([65c9b9f](https://github.com/perfect-panel/ppanel-web/commit/65c9b9f))

### ✨ Features

- **api**: Add getClient API endpoint to retrieve subscription applications ([7a279e6](https://github.com/perfect-panel/ppanel-web/commit/7a279e6))
- **marketing**: Add marketing management features and localization updates ([ea08de0](https://github.com/perfect-panel/ppanel-web/commit/ea08de0))
- **protocol**: Add template preview functionality with localization support ([0448d21](https://github.com/perfect-panel/ppanel-web/commit/0448d21))
- **subscribe**: Update subscription management localization and add new fields ([1d9b0a4](https://github.com/perfect-panel/ppanel-web/commit/1d9b0a4))

### 🐛 Bug Fixes

- **bun**: Update bun.lockb to reflect dependency changes ([bbcd018](https://github.com/perfect-panel/ppanel-web/commit/bbcd018))
- **editor**: Add Go template editor component and update related schemas ([9d9c3cd](https://github.com/perfect-panel/ppanel-web/commit/9d9c3cd))
- **editor**: Enhance Go Template Editor to support trimmed template tags and improve range/end matching ([641ed5e](https://github.com/perfect-panel/ppanel-web/commit/641ed5e))
- **editor**: Enhance Go Template Editor with schema support and improved completion ([5b21d8a](https://github.com/perfect-panel/ppanel-web/commit/5b21d8a))
- **enhaced-input**: Disable autocomplete for EnhancedInput component ([f190c68](https://github.com/perfect-panel/ppanel-web/commit/f190c68))
- **global**: Add user agent limit settings to subscription configuration ([822416d](https://github.com/perfect-panel/ppanel-web/commit/822416d))
- **locales**: Update 'conf' output format to use uppercase 'CONF' ([fce9119](https://github.com/perfect-panel/ppanel-web/commit/fce9119))
- **locales**: Update "userAccount" label to "user" in multiple localization files ([48415e9](https://github.com/perfect-panel/ppanel-web/commit/48415e9))
- **protocol-form**: Swap 'description' and 'user_agent' columns for improved clarity in the table ([72a4106](https://github.com/perfect-panel/ppanel-web/commit/72a4106))
- **protocol-form**: Update protocol options descriptions for clarity and add new security and transport options ([e5d4deb](https://github.com/perfect-panel/ppanel-web/commit/e5d4deb))
- **protocol**: Add 'conf' output format option and update translations ([292efdf](https://github.com/perfect-panel/ppanel-web/commit/292efdf))
- **protpcp-form**: Rename 'schema' to 'scheme' for consistency across the application ([6ab2ba9](https://github.com/perfect-panel/ppanel-web/commit/6ab2ba9))
- **register**: Update localization files to include trial subscription settings and descriptions ([33daa1f](https://github.com/perfect-panel/ppanel-web/commit/33daa1f))
- **system**: Add time unit translations for user registration settings in multiple languages ([296a6c1](https://github.com/perfect-panel/ppanel-web/commit/296a6c1))
- Update privacy policy and terms of service schemas to use correct field names ([0e6ba5b](https://github.com/perfect-panel/ppanel-web/commit/0e6ba5b))

<a name="readme-top"></a>

# Changelog

# [1.2.0](https://github.com/perfect-panel/ppanel-web/compare/v1.1.5...v1.2.0) (2025-08-04)

### ♻ Code Refactoring

- **view**: System and Auth Control ([b2b4a95](https://github.com/perfect-panel/ppanel-web/commit/b2b4a95))

### ✨ Features

- **netlify**: Add Netlify configuration for admin and user apps with Next.js plugin ([b4d4f59](https://github.com/perfect-panel/ppanel-web/commit/b4d4f59))

<a name="readme-top"></a>

# Changelog

## [1.1.5](https://github.com/perfect-panel/ppanel-web/compare/v1.1.4...v1.1.5) (2025-07-26)

### 🐛 Bug Fixes

- **subscribe**: Filter out items that are not marked as visible in subscription list ([32253e3](https://github.com/perfect-panel/ppanel-web/commit/32253e3))

## [1.1.4](https://github.com/perfect-panel/ppanel-web/compare/v1.1.3...v1.1.4) (2025-07-25)

### 🐛 Bug Fixes

- **locales**: Simplify "show" label in subscription localization files ([d53a006](https://github.com/perfect-panel/ppanel-web/commit/d53a006))
- **order**: Preserve last successful order on error during order creation ([2fb98be](https://github.com/perfect-panel/ppanel-web/commit/2fb98be))
- **subscribe**: Filter out hidden items in subscription list display ([634be37](https://github.com/perfect-panel/ppanel-web/commit/634be37))

## [1.1.3](https://github.com/perfect-panel/ppanel-web/compare/v1.1.2...v1.1.3) (2025-07-24)

### 🐛 Bug Fixes

- **auth**: Implement user redirection to dashboard upon authentication ([f84f98c](https://github.com/perfect-panel/ppanel-web/commit/f84f98c))

## [1.1.2](https://github.com/perfect-panel/ppanel-web/compare/v1.1.1...v1.1.2) (2025-07-24)

### 🐛 Bug Fixes

- **billing**: Add display for gift amount in subscription billing ([04af2f9](https://github.com/perfect-panel/ppanel-web/commit/04af2f9))
- **order**: Update subscription cell to display name and quantity ([96eba17](https://github.com/perfect-panel/ppanel-web/commit/96eba17))
- **tool**: Added API for obtaining version, updated version information display logic ([2675034](https://github.com/perfect-panel/ppanel-web/commit/2675034))

<a name="readme-top"></a>

# Changelog

## [1.1.1](https://github.com/perfect-panel/ppanel-web/compare/v1.1.0...v1.1.1) (2025-07-20)

### 🐛 Bug Fixes

- **node-table**: Update translations for headers and no data display ([eec0b12](https://github.com/perfect-panel/ppanel-web/commit/eec0b12))
- **rules**: Change rule type from 'auto' to 'default' and update ([3e290d7](https://github.com/perfect-panel/ppanel-web/commit/3e290d7))
- **rules**: Update rule settings ([3304a55](https://github.com/perfect-panel/ppanel-web/commit/3304a55))
- **subscribe-form**: Optimize discount calculation logic and debounce updates ([166e48f](https://github.com/perfect-panel/ppanel-web/commit/166e48f))
- **tutorial**: Comment out unused getVersion function and simplify getVersionPath ([7cdc6bd](https://github.com/perfect-panel/ppanel-web/commit/7cdc6bd))
- **tutorial**: Return latest version in case of fetch error ([1fb305e](https://github.com/perfect-panel/ppanel-web/commit/1fb305e))

<a name="readme-top"></a>

# Changelog

# [1.1.0](https://github.com/perfect-panel/ppanel-web/compare/v1.0.2...v1.1.0) (2025-07-06)

### ✨ Features

- **view**: Add AnyTLS protocol support and enhance node configuration options ([bcfb10a](https://github.com/perfect-panel/ppanel-web/commit/bcfb10a))

<a name="readme-top"></a>

# Changelog

## [1.0.2](https://github.com/perfect-panel/ppanel-web/compare/v1.0.1...v1.0.2) (2025-06-29)

### 🐛 Bug Fixes

- **subscription**: User subscription information ([7d18ff6](https://github.com/perfect-panel/ppanel-web/commit/7d18ff6))

## [1.0.1](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0...v1.0.1) (2025-04-28)

### 🐛 Bug Fixes

- **payment**: Disable webhook_secret field in PaymentForm component ([d323af8](https://github.com/perfect-panel/ppanel-web/commit/d323af8))
- **recharge**: Set balance prop to false in PaymentMethods component ([356ae5b](https://github.com/perfect-panel/ppanel-web/commit/356ae5b))

# 1.0.0 (2025-04-24)

### ♻ Code Refactoring

- **api**: Sort and Announcement ([38d5616](https://github.com/perfect-panel/ppanel-web/commit/38d5616))
- **auth**: Refactor user authorization handling and improve error logging ([68bc18f](https://github.com/perfect-panel/ppanel-web/commit/68bc18f))
- **config**: GenerateMetadata ([a0bb101](https://github.com/perfect-panel/ppanel-web/commit/a0bb101))
- **config**: Simplify environment variable handling and improve build script ([cf54d0f](https://github.com/perfect-panel/ppanel-web/commit/cf54d0f))
- **config**: Viewport ([24b8601](https://github.com/perfect-panel/ppanel-web/commit/24b8601))
- **core**: Restructure project for better module separation ([9d0cb8b](https://github.com/perfect-panel/ppanel-web/commit/9d0cb8b))
- **deps**: Update ([19837a1](https://github.com/perfect-panel/ppanel-web/commit/19837a1))
- **empty**: Content ([aa4c667](https://github.com/perfect-panel/ppanel-web/commit/aa4c667))
- **payment**: Reconstruct the payment page ([7109472](https://github.com/perfect-panel/ppanel-web/commit/7109472))
- **sbscribe**: Rename and reorganize components for better structure and clarity ([5e5e4ed](https://github.com/perfect-panel/ppanel-web/commit/5e5e4ed))
- **ui**: Dependencies ([727d779](https://github.com/perfect-panel/ppanel-web/commit/727d779))
- **ui**: Layout ([9262d7d](https://github.com/perfect-panel/ppanel-web/commit/9262d7d))
- **ui**: Optimize document display ([2ca2992](https://github.com/perfect-panel/ppanel-web/commit/2ca2992))
- Enhance user navigation dropdown ui and styling ([d2732e6](https://github.com/perfect-panel/ppanel-web/commit/d2732e6))
- Reduce code complexity and improve readability ([e11f18c](https://github.com/perfect-panel/ppanel-web/commit/e11f18c))

### ⚡ Performance Improvements

- **subscribe**: Form discount price ([059a892](https://github.com/perfect-panel/ppanel-web/commit/059a892))

### ✨ Features

- **accounts**: Update third-party account binding and unbinding ([1841552](https://github.com/perfect-panel/ppanel-web/commit/1841552))
- **ad**: Advertise ([b1105cd](https://github.com/perfect-panel/ppanel-web/commit/b1105cd))
- **admin**: Add application and rule management entries to localization files ([8b43e69](https://github.com/perfect-panel/ppanel-web/commit/8b43e69))
- **affiliate**: Add Affiliate component with commission display and invite link functionality ([4aea4e8](https://github.com/perfect-panel/ppanel-web/commit/4aea4e8))
- **affiliate**: Affiliate Detail ([a782c17](https://github.com/perfect-panel/ppanel-web/commit/a782c17))
- **affiliate**: Commission Rate ([5eec430](https://github.com/perfect-panel/ppanel-web/commit/5eec430))
- **affiliate**: Update affiliate component to display total commission and improve data fetching ([cc834ca](https://github.com/perfect-panel/ppanel-web/commit/cc834ca))
- **announcement**: Popup and pinned ([f3680a7](https://github.com/perfect-panel/ppanel-web/commit/f3680a7))
- **api**: Add an interface to obtain user subscription details, update related type definitions and localized text ([cf5c39c](https://github.com/perfect-panel/ppanel-web/commit/cf5c39c))
- **api**: Add CheckoutOrder request and response types, and update user purchase request parameters ([48a1b97](https://github.com/perfect-panel/ppanel-web/commit/48a1b97))
- **api**: Add CheckoutOrder request and response types, and update user purchase request parameters ([dddc21c](https://github.com/perfect-panel/ppanel-web/commit/dddc21c))
- **api**: Add new subscription properties and locale support for deduction ratios and reset cycles ([fec80f5](https://github.com/perfect-panel/ppanel-web/commit/fec80f5))
- **api**: Add Time Period Configuration ([837157c](https://github.com/perfect-panel/ppanel-web/commit/837157c))
- **api**: Telegram ([17ce96a](https://github.com/perfect-panel/ppanel-web/commit/17ce96a))
- **auth-control**: Adding phone number labels to mobile verification configurations in multiple languages ([046740f](https://github.com/perfect-panel/ppanel-web/commit/046740f))
- **auth-control**: Update general ([3883646](https://github.com/perfect-panel/ppanel-web/commit/3883646))
- **auth**: Add email and SMS code sending functionality with localization updates ([57eaa55](https://github.com/perfect-panel/ppanel-web/commit/57eaa55))
- **auth**: Add Oauth configuration for Telegram, Facebook, Google, Github, and Apple ([18ee600](https://github.com/perfect-panel/ppanel-web/commit/18ee600))
- **auth**: Add privacy policy link to the footer ([8e16ef1](https://github.com/perfect-panel/ppanel-web/commit/8e16ef1))
- **auth**: Add SMS and email configuration options to global store and update localization ([4acf7b1](https://github.com/perfect-panel/ppanel-web/commit/4acf7b1))
- **auth**: Add type parameter to SendCode and update related API typings ([4198871](https://github.com/perfect-panel/ppanel-web/commit/4198871))
- **auth**: Enhance user registration with invite handling and logo display ([207bc24](https://github.com/perfect-panel/ppanel-web/commit/207bc24))
- **auth**: Redirect user after OAuth login and add logos icon collection ([aa6dda8](https://github.com/perfect-panel/ppanel-web/commit/aa6dda8))
- **auth**: Refactor mobile authentication config to support whitelist functionality ([c761ec7](https://github.com/perfect-panel/ppanel-web/commit/c761ec7))
- **billing**: Update Billing ([078fc9d](https://github.com/perfect-panel/ppanel-web/commit/078fc9d))
- **cdn**: Add CDN URL configuration and update related references ([0c90733](https://github.com/perfect-panel/ppanel-web/commit/0c90733))
- **config**: Add application selection and encryption settings to configuration form ([88b3504](https://github.com/perfect-panel/ppanel-web/commit/88b3504))
- **config**: FormatBytes ([9251a09](https://github.com/perfect-panel/ppanel-web/commit/9251a09))
- **config**: Protocol type ([a3b45b4](https://github.com/perfect-panel/ppanel-web/commit/a3b45b4))
- **config**: Update encryption fields in configuration form and refactor OAuth callback parameters ([652e032](https://github.com/perfect-panel/ppanel-web/commit/652e032))
- **config**: Webhook Domain ([01e06c6](https://github.com/perfect-panel/ppanel-web/commit/01e06c6))
- **dashboard**: Optimization ([5b3f4b4](https://github.com/perfect-panel/ppanel-web/commit/5b3f4b4))
- **dashboard**: Statistics ([2926abc](https://github.com/perfect-panel/ppanel-web/commit/2926abc))
- **device**: Modify IMEI to device identifier support ([e3f9ef6](https://github.com/perfect-panel/ppanel-web/commit/e3f9ef6))
- **email**: Add traffic exhaustion template ([bb3bd7b](https://github.com/perfect-panel/ppanel-web/commit/bb3bd7b))
- **favicon**: Update SVG favicon design for admin and user interfaces ([1d91738](https://github.com/perfect-panel/ppanel-web/commit/1d91738))
- **formatting**: Update differenceInDays function to return whole days or two decimal places ([bf58f25](https://github.com/perfect-panel/ppanel-web/commit/bf58f25))
- **form**: Make version field optional and set default value; update site domain placeholder for clarity ([42ba9e8](https://github.com/perfect-panel/ppanel-web/commit/42ba9e8))
- **global**: Add custom data ([6dbebd1](https://github.com/perfect-panel/ppanel-web/commit/6dbebd1))
- **global**: Add SMS configuration options to global store ([39a9ce6](https://github.com/perfect-panel/ppanel-web/commit/39a9ce6))
- **header**: Update locales ([bfb6c27](https://github.com/perfect-panel/ppanel-web/commit/bfb6c27))
- **imei**: Add IMEI related internationalization support and menu items ([13c3337](https://github.com/perfect-panel/ppanel-web/commit/13c3337))
- **input**: Add minimum value constraint and enhance number handling in EnhancedInput ([ce31972](https://github.com/perfect-panel/ppanel-web/commit/ce31972))
- **input**: Add minimum value constraint and enhance number handling in EnhancedInput ([94822d9](https://github.com/perfect-panel/ppanel-web/commit/94822d9))
- **loading**: Add loading components and integrate them in Providers ([d5847fa](https://github.com/perfect-panel/ppanel-web/commit/d5847fa))
- **loading**: Replace loading animation with a simpler spinner and loading text ([f72df3a](https://github.com/perfect-panel/ppanel-web/commit/f72df3a))
- **loading**: Replace loading animation with a simpler spinner and loading text ([b8316bb](https://github.com/perfect-panel/ppanel-web/commit/b8316bb))
- **locale**: Add Persian ([93a0a88](https://github.com/perfect-panel/ppanel-web/commit/93a0a88))
- **locales**: Add area code and telephone fields to user forms in multiple languages ([9b8258c](https://github.com/perfect-panel/ppanel-web/commit/9b8258c))
- **locales**: Add description information of communication keys and encryption methods to enhance client configuration capabilities ([d1f5a9b](https://github.com/perfect-panel/ppanel-web/commit/d1f5a9b))
- **locales**: Add kick offline confirmation and success messages in multiple languages ([5db5343](https://github.com/perfect-panel/ppanel-web/commit/5db5343))
- **locales**: Add multiple languages ([b243ab9](https://github.com/perfect-panel/ppanel-web/commit/b243ab9))
- **locales**: Replace 'nodeGroupId' with 'groupId' in multiple language files for consistency ([a4e9d5d](https://github.com/perfect-panel/ppanel-web/commit/a4e9d5d))
- **locales**: Update 'deductBalance' to 'giftAmount' across multiple languages and fix newline in announcement.json ([70497af](https://github.com/perfect-panel/ppanel-web/commit/70497af))
- **locales**: Update 'sms' to 'mobile' in authentication methods across multiple languages ([fea2171](https://github.com/perfect-panel/ppanel-web/commit/fea2171))
- **log**: Add message log retrieval functionality and update related typings ([1c0ecae](https://github.com/perfect-panel/ppanel-web/commit/1c0ecae))
- **node-form**: Update number input fields to enforce step, min, and max values ([3f7b6d1](https://github.com/perfect-panel/ppanel-web/commit/3f7b6d1))
- **node-subscription**: Add copy functionality for columns ([3a81e37](https://github.com/perfect-panel/ppanel-web/commit/3a81e37))
- **node**: Add NodeStatus ([c712624](https://github.com/perfect-panel/ppanel-web/commit/c712624))
- **node**: Add protocol ([301b635](https://github.com/perfect-panel/ppanel-web/commit/301b635))
- **node**: Add serverKey ([25ce37e](https://github.com/perfect-panel/ppanel-web/commit/25ce37e))
- **node**: Add status ([c06372b](https://github.com/perfect-panel/ppanel-web/commit/c06372b))
- **node**: Add tags ([f408fdf](https://github.com/perfect-panel/ppanel-web/commit/f408fdf))
- **node**: Move the node configuration to the server module ([7f0f5ce](https://github.com/perfect-panel/ppanel-web/commit/7f0f5ce))
- **oauth**: Add certification component for handling OAuth login callbacks and improve user authentication flow ([5ed04c0](https://github.com/perfect-panel/ppanel-web/commit/5ed04c0))
- **oauth**: Implement OAuth token retrieval and refactor login callback handling ([40a6f7c](https://github.com/perfect-panel/ppanel-web/commit/40a6f7c))
- **oauth**: Refactor platform parameter handling and improve logout redirection logic ([8346c85](https://github.com/perfect-panel/ppanel-web/commit/8346c85))
- **oauth**: Update OAuth login handling to use callback parameter and improve URL parameter retrieval ([9227411](https://github.com/perfect-panel/ppanel-web/commit/9227411))
- **payment**: Add bank card payment ([7fa3a57](https://github.com/perfect-panel/ppanel-web/commit/7fa3a57))
- **payment**: Add isEdit prop to PaymentForm and disable fields when editing ([85f55de](https://github.com/perfect-panel/ppanel-web/commit/85f55de))
- **platform**: Update platform naming and add keywords and custom HTML fields ([6384237](https://github.com/perfect-panel/ppanel-web/commit/6384237))
- **privacy-policy**: Add privacy policy related text and links ([baa68f0](https://github.com/perfect-panel/ppanel-web/commit/baa68f0))
- **profile**: Update localization strings and enhance third-party account binding ([2d1effb](https://github.com/perfect-panel/ppanel-web/commit/2d1effb))
- **relay**: Add relay mode configuration and update related schemas ([3cc9477](https://github.com/perfect-panel/ppanel-web/commit/3cc9477))
- **release**: Extend supported platforms for Docker images, closes [#9](https://github.com/perfect-panel/ppanel-web/issues/9) ([e3a31eb](https://github.com/perfect-panel/ppanel-web/commit/e3a31eb))
- **schema**: Add security field to hysteria2 and tuic schemas ([cd59d44](https://github.com/perfect-panel/ppanel-web/commit/cd59d44))
- **site**: Added localization support for custom HTML and keyword fields ([f9d7736](https://github.com/perfect-panel/ppanel-web/commit/f9d7736))
- **sms**: Update locales ([938363b](https://github.com/perfect-panel/ppanel-web/commit/938363b))
- **stats**: Replace dynamic stat fetching with environment constants for user, server, and location counts ([46ae166](https://github.com/perfect-panel/ppanel-web/commit/46ae166))
- **subscribe**: Add 'sold' column to SubscribeTable and update inventory terminology ([19619fd](https://github.com/perfect-panel/ppanel-web/commit/19619fd))
- **subscribe**: Add reset_time to API typings and update unsubscribe logic ([eeea165](https://github.com/perfect-panel/ppanel-web/commit/eeea165))
- **subscribe**: Add subscribe_discount type ([f99c604](https://github.com/perfect-panel/ppanel-web/commit/f99c604))
- **subscribe**: Add subscription credits ([5bc7905](https://github.com/perfect-panel/ppanel-web/commit/5bc7905))
- **subscribe**: Add unit time ([39d07ec](https://github.com/perfect-panel/ppanel-web/commit/39d07ec))
- **subscribe**: Add unsubscribe functionality with confirmation messages and localized strings ([b2a2f42](https://github.com/perfect-panel/ppanel-web/commit/b2a2f42))
- **subscribe**: Improve error handling in subscription forms and update component props ([d28a10b](https://github.com/perfect-panel/ppanel-web/commit/d28a10b))
- **subscribe**: Improve layout and styling in subscription components ([5766376](https://github.com/perfect-panel/ppanel-web/commit/5766376))
- **subscribe**: Move subscription configuration and application to subscription module ([f90d4d2](https://github.com/perfect-panel/ppanel-web/commit/f90d4d2))
- **subscribe**: Update SubscribeTable component to use API.SubscribeItem type and ensure proper type casting ([f26f1c2](https://github.com/perfect-panel/ppanel-web/commit/f26f1c2))
- **subscribe**: Update suffix from 'MB' to 'Mbps' and enhance speed limit display logic ([3547bb1](https://github.com/perfect-panel/ppanel-web/commit/3547bb1))
- **subscription**: Improve layout and organization of subscription detail tabs ([e4630f8](https://github.com/perfect-panel/ppanel-web/commit/e4630f8))
- **subscription**: Add delete user subscription functionality ([1fc3a10](https://github.com/perfect-panel/ppanel-web/commit/1fc3a10))
- **subscription**: Add localized messages for existing subscriptions and deletion restrictions ([e8a72d5](https://github.com/perfect-panel/ppanel-web/commit/e8a72d5))
- **subscription**: Refactor subscription handling and update imports for better organization ([2215c7f](https://github.com/perfect-panel/ppanel-web/commit/2215c7f))
- **table**: Add sorting support for Node and subscription columns ([27924b0](https://github.com/perfect-panel/ppanel-web/commit/27924b0))
- **table**: Supports drag and drop sorting ([2f56ef5](https://github.com/perfect-panel/ppanel-web/commit/2f56ef5))
- **timeline**: Simplify timeline component layout and remove commented-out code ([fbad3b0](https://github.com/perfect-panel/ppanel-web/commit/fbad3b0))
- **tos**: Display data ([6024454](https://github.com/perfect-panel/ppanel-web/commit/6024454))
- **tutorial**: Add common tutorial list ([872252c](https://github.com/perfect-panel/ppanel-web/commit/872252c))
- **tutorial**: Fetch the latest tutorial version from GitHub API for dynamic URL generation ([28f8c78](https://github.com/perfect-panel/ppanel-web/commit/28f8c78))
- **ui**: System Tool ([1836980](https://github.com/perfect-panel/ppanel-web/commit/1836980))
- **ui**: Update homepage data ([8425b13](https://github.com/perfect-panel/ppanel-web/commit/8425b13))
- **ui**: Update input components and enhance card minimum width for better layout ([8a02310](https://github.com/perfect-panel/ppanel-web/commit/8a02310))
- **user**: Add 'gift_amount' field and update related references in user services and components ([b13c77e](https://github.com/perfect-panel/ppanel-web/commit/b13c77e))
- **user**: Add telephone input with area code selection and update localization ([585b99c](https://github.com/perfect-panel/ppanel-web/commit/585b99c))
- **user**: Add user Detail ([3a3d223](https://github.com/perfect-panel/ppanel-web/commit/3a3d223))
- **user**: Add User Detail ([fdaf11b](https://github.com/perfect-panel/ppanel-web/commit/fdaf11b))
- **user**: Integrate subscription list into user management, update request parameters and types ([8d49dac](https://github.com/perfect-panel/ppanel-web/commit/8d49dac))
- Update Auth Control ([c59742a](https://github.com/perfect-panel/ppanel-web/commit/c59742a))

### 🎫 Chores

- **config**: Entry locale ([5737331](https://github.com/perfect-panel/ppanel-web/commit/5737331))
- **deps**: Update package dependencies across multiple projects for improved stability and performance ([b01a5bc](https://github.com/perfect-panel/ppanel-web/commit/b01a5bc))
- **init**: Project initialization ([829edfa](https://github.com/perfect-panel/ppanel-web/commit/829edfa))
- **merge**: Add advertising module and device settings ([0130e02](https://github.com/perfect-panel/ppanel-web/commit/0130e02))
- **merge**: Bump version to 1.0.0-beta.26 and update changelog ([3222016](https://github.com/perfect-panel/ppanel-web/commit/3222016))
- **release**: V1.0.0-beta.1 [skip ci] ([7284d1c](https://github.com/perfect-panel/ppanel-web/commit/7284d1c))
- **release**: V1.0.0-beta.10 [skip ci] ([5cf573a](https://github.com/perfect-panel/ppanel-web/commit/5cf573a))
- **release**: V1.0.0-beta.11 [skip ci] ([1f29506](https://github.com/perfect-panel/ppanel-web/commit/1f29506))
- **release**: V1.0.0-beta.12 [skip ci] ([4418c47](https://github.com/perfect-panel/ppanel-web/commit/4418c47))
- **release**: V1.0.0-beta.13 [skip ci] ([23c974a](https://github.com/perfect-panel/ppanel-web/commit/23c974a))
- **release**: V1.0.0-beta.14 [skip ci] ([0fb0d8b](https://github.com/perfect-panel/ppanel-web/commit/0fb0d8b))
- **release**: V1.0.0-beta.15 [skip ci] ([b2e8fad](https://github.com/perfect-panel/ppanel-web/commit/b2e8fad))
- **release**: V1.0.0-beta.16 [skip ci] ([c3eff0a](https://github.com/perfect-panel/ppanel-web/commit/c3eff0a))
- **release**: V1.0.0-beta.17 [skip ci] ([5b64389](https://github.com/perfect-panel/ppanel-web/commit/5b64389))
- **release**: V1.0.0-beta.18 [skip ci] ([4a00233](https://github.com/perfect-panel/ppanel-web/commit/4a00233))
- **release**: V1.0.0-beta.19 [skip ci] ([0f15fb8](https://github.com/perfect-panel/ppanel-web/commit/0f15fb8))
- **release**: V1.0.0-beta.2 [skip ci] ([087c36c](https://github.com/perfect-panel/ppanel-web/commit/087c36c))
- **release**: V1.0.0-beta.20 [skip ci] ([bbd44f0](https://github.com/perfect-panel/ppanel-web/commit/bbd44f0))
- **release**: V1.0.0-beta.21 [skip ci] ([ca642c2](https://github.com/perfect-panel/ppanel-web/commit/ca642c2))
- **release**: V1.0.0-beta.22 [skip ci] ([c0fb34f](https://github.com/perfect-panel/ppanel-web/commit/c0fb34f))
- **release**: V1.0.0-beta.23 [skip ci] ([cf1d66d](https://github.com/perfect-panel/ppanel-web/commit/cf1d66d))
- **release**: V1.0.0-beta.24 [skip ci] ([01a3aa0](https://github.com/perfect-panel/ppanel-web/commit/01a3aa0))
- **release**: V1.0.0-beta.25 [skip ci] ([047a698](https://github.com/perfect-panel/ppanel-web/commit/047a698))
- **release**: V1.0.0-beta.26 [skip ci] ([79edea7](https://github.com/perfect-panel/ppanel-web/commit/79edea7))
- **release**: V1.0.0-beta.27 [skip ci] ([092477b](https://github.com/perfect-panel/ppanel-web/commit/092477b))
- **release**: V1.0.0-beta.27 [skip ci] ([85fdc36](https://github.com/perfect-panel/ppanel-web/commit/85fdc36))
- **release**: V1.0.0-beta.28 [skip ci] ([786ba0e](https://github.com/perfect-panel/ppanel-web/commit/786ba0e))
- **release**: V1.0.0-beta.28 [skip ci] ([d10ecc9](https://github.com/perfect-panel/ppanel-web/commit/d10ecc9))
- **release**: V1.0.0-beta.29 [skip ci] ([29bc3c7](https://github.com/perfect-panel/ppanel-web/commit/29bc3c7))
- **release**: V1.0.0-beta.3 [skip ci] ([cd49427](https://github.com/perfect-panel/ppanel-web/commit/cd49427))
- **release**: V1.0.0-beta.30 [skip ci] ([db0d9e0](https://github.com/perfect-panel/ppanel-web/commit/db0d9e0))
- **release**: V1.0.0-beta.31 [skip ci] ([aa1d426](https://github.com/perfect-panel/ppanel-web/commit/aa1d426))
- **release**: V1.0.0-beta.32 [skip ci] ([fa56eb8](https://github.com/perfect-panel/ppanel-web/commit/fa56eb8))
- **release**: V1.0.0-beta.33 [skip ci] ([383638f](https://github.com/perfect-panel/ppanel-web/commit/383638f))
- **release**: V1.0.0-beta.34 [skip ci] ([7023875](https://github.com/perfect-panel/ppanel-web/commit/7023875))
- **release**: V1.0.0-beta.4 [skip ci] ([10d322f](https://github.com/perfect-panel/ppanel-web/commit/10d322f))
- **release**: V1.0.0-beta.5 [skip ci] ([f275f01](https://github.com/perfect-panel/ppanel-web/commit/f275f01))
- **release**: V1.0.0-beta.6 [skip ci] ([dd23279](https://github.com/perfect-panel/ppanel-web/commit/dd23279))
- **release**: V1.0.0-beta.7 [skip ci] ([f60d40c](https://github.com/perfect-panel/ppanel-web/commit/f60d40c))
- **release**: V1.0.0-beta.8 [skip ci], closes [#9](https://github.com/perfect-panel/ppanel-web/issues/9) ([a593eac](https://github.com/perfect-panel/ppanel-web/commit/a593eac))
- **release**: V1.0.0-beta.9 [skip ci] ([855d1b0](https://github.com/perfect-panel/ppanel-web/commit/855d1b0))
- **ui**: Update package dependencies for improved stability and performance ([25da429](https://github.com/perfect-panel/ppanel-web/commit/25da429))
- Merge branch 'beta' into develop ([f219c52](https://github.com/perfect-panel/ppanel-web/commit/f219c52))
- Update changelog, enhance prepare script, and add openapi command ([a93db4e](https://github.com/perfect-panel/ppanel-web/commit/a93db4e))

### 🐛 Bug Fixes

- **admin**: Hidden versions and system upgrades ([64cd842](https://github.com/perfect-panel/ppanel-web/commit/64cd842))
- **admin**: Modify the label type in the rule form to a string array ([a7aa5fe](https://github.com/perfect-panel/ppanel-web/commit/a7aa5fe))
- **affiliate**: Update user identifier ([35f92c9](https://github.com/perfect-panel/ppanel-web/commit/35f92c9))
- **api**: Fix type error in API request and add return URL parameter ([ee286dd](https://github.com/perfect-panel/ppanel-web/commit/ee286dd))
- **api**: PreCreateOrder ([ca747f5](https://github.com/perfect-panel/ppanel-web/commit/ca747f5))
- **api**: Purge ([98c1c30](https://github.com/perfect-panel/ppanel-web/commit/98c1c30))
- **api**: Remove redundant requestType parameter in appleLoginCallback ([0aa5d5b](https://github.com/perfect-panel/ppanel-web/commit/0aa5d5b))
- **api**: Rename app-related functions and types to application for consistency ([9d8b814](https://github.com/perfect-panel/ppanel-web/commit/9d8b814))
- **api**: Replace 'deduction' with 'gift_amount' and add 'commission' field in type definitions ([77edf1d](https://github.com/perfect-panel/ppanel-web/commit/77edf1d))
- **api**: Server and order ([255bd82](https://github.com/perfect-panel/ppanel-web/commit/255bd82))
- **api**: Statistics ([7962162](https://github.com/perfect-panel/ppanel-web/commit/7962162))
- **api**: Subscribe token ([1932ba7](https://github.com/perfect-panel/ppanel-web/commit/1932ba7))
- **api**: Update API type definitions to replace 'deduction' with 'gift_amount' and make 'commission' field optional ([c2af060](https://github.com/perfect-panel/ppanel-web/commit/c2af060))
- **api**: Update Model ([39aaa73](https://github.com/perfect-panel/ppanel-web/commit/39aaa73))
- **api**: Update subscription_protocol to subscribe_type for consistency across services ([b6da51b](https://github.com/perfect-panel/ppanel-web/commit/b6da51b))
- **auth-control**: Fix citation error for platform values ([c940f3c](https://github.com/perfect-panel/ppanel-web/commit/c940f3c))
- **auth-control**: Fix citation error for platform values ([28813d2](https://github.com/perfect-panel/ppanel-web/commit/28813d2))
- **auth-control**: Rename phone_variable to phone_number in mobile verification configuration ([e5455aa](https://github.com/perfect-panel/ppanel-web/commit/e5455aa))
- **auth**: Add error handling to form submission and reset Turnstile on failure ([715d011](https://github.com/perfect-panel/ppanel-web/commit/715d011))
- **auth**: Change Textarea value to defaultValue for client_secret in Apple auth page ([69fc670](https://github.com/perfect-panel/ppanel-web/commit/69fc670))
- **auth**: Refactor forms to use Turnstile ref for reset functionality ([320a7dc](https://github.com/perfect-panel/ppanel-web/commit/320a7dc))
- **auth**: Refactor reset password form to simplify code input and update placeholder text ([23833b4](https://github.com/perfect-panel/ppanel-web/commit/23833b4))
- **auth**: Refactor user authentication forms to remove global store dependency and improve type handling ([12026b0](https://github.com/perfect-panel/ppanel-web/commit/12026b0))
- **auth**: Remove unused telephone code login function and update typings for telephone login requests ([7239685](https://github.com/perfect-panel/ppanel-web/commit/7239685))
- **auth**: Require minimum length for invite string when forced invite is enabled ([a604f28](https://github.com/perfect-panel/ppanel-web/commit/a604f28))
- **auth**: Simplify email verification code input rendering ([6f7bc37](https://github.com/perfect-panel/ppanel-web/commit/6f7bc37))
- **auth**: Update authentication configuration and localization strings ([47f2c58](https://github.com/perfect-panel/ppanel-web/commit/47f2c58))
- **auth**: Update email verification logic to use domain suffix check ([62662bb](https://github.com/perfect-panel/ppanel-web/commit/62662bb))
- **auth**: Update user authentication flow to include email and phone code verification ([5d078fd](https://github.com/perfect-panel/ppanel-web/commit/5d078fd))
- **auth**: Update UserCheckForm to use setInitialValues and modify onSubmit type ([c984c0d](https://github.com/perfect-panel/ppanel-web/commit/c984c0d))
- **billing**: ExpiryDate ([e85e545](https://github.com/perfect-panel/ppanel-web/commit/e85e545))
- **billing**: I18n and styles ([81e0f21](https://github.com/perfect-panel/ppanel-web/commit/81e0f21))
- **changelog**: Update change log style ([cfa3fc0](https://github.com/perfect-panel/ppanel-web/commit/cfa3fc0))
- **config**: AlipayF2F ([6c07107](https://github.com/perfect-panel/ppanel-web/commit/6c07107))
- **config**: Bugs ([f57e40c](https://github.com/perfect-panel/ppanel-web/commit/f57e40c))
- **config**: Checkout Order ([a31e763](https://github.com/perfect-panel/ppanel-web/commit/a31e763))
- **config**: FormatBytes ([bbc2da0](https://github.com/perfect-panel/ppanel-web/commit/bbc2da0))
- **config**: NoStore ([2cc18cf](https://github.com/perfect-panel/ppanel-web/commit/2cc18cf))
- **config**: Runtime env ([a1e4999](https://github.com/perfect-panel/ppanel-web/commit/a1e4999))
- **config**: Status Percentag ([8f322fb](https://github.com/perfect-panel/ppanel-web/commit/8f322fb))
- **config**: SubLink ([1c61966](https://github.com/perfect-panel/ppanel-web/commit/1c61966))
- **config**: Subscribe Link ([11ea821](https://github.com/perfect-panel/ppanel-web/commit/11ea821))
- **content**: Parse subscription description and display features with icons ([3c5542a](https://github.com/perfect-panel/ppanel-web/commit/3c5542a))
- **controller**: Order status ([8c6a097](https://github.com/perfect-panel/ppanel-web/commit/8c6a097))
- **coupon**: Rename 'server' field to 'subscribe' in coupon form and update coupon update request type ([f8b6d82](https://github.com/perfect-panel/ppanel-web/commit/f8b6d82))
- **dashboard**: Improve URL encoding for subscription links and enhance success message handling ([4983c33](https://github.com/perfect-panel/ppanel-web/commit/4983c33))
- **dashboard**: Update date display to use start_time if available ([e551232](https://github.com/perfect-panel/ppanel-web/commit/e551232))
- **dashboard**: Correct progress value calculations and update groupId accessor ([36c7667](https://github.com/perfect-panel/ppanel-web/commit/36c7667))
- **dashboard**: Display subscription creation date in user dashboard ([d0e6df0](https://github.com/perfect-panel/ppanel-web/commit/d0e6df0))
- **dashboard**: Format Bytes ([d8b0bd9](https://github.com/perfect-panel/ppanel-web/commit/d8b0bd9))
- **dashboard**: Update icon imports for platform consistency and adjust icon size ([3e8912e](https://github.com/perfect-panel/ppanel-web/commit/3e8912e))
- **dashboard**: Update platform detection logic and improve layout responsiveness ([b0aa364](https://github.com/perfect-panel/ppanel-web/commit/b0aa364))
- **deps**: Remove outdated @iconify/react dependency and add iconify-json packages ([d6fbc38](https://github.com/perfect-panel/ppanel-web/commit/d6fbc38))
- **deps**: Typescript config ([34e24b8](https://github.com/perfect-panel/ppanel-web/commit/34e24b8))
- **deps**: Update clipboard ([5572710](https://github.com/perfect-panel/ppanel-web/commit/5572710))
- **editor**: Change value ([4fdfeb2](https://github.com/perfect-panel/ppanel-web/commit/4fdfeb2))
- **email**: Update platform configuration handling to use current ref for consistency ([c90175b](https://github.com/perfect-panel/ppanel-web/commit/c90175b))
- **footer**: Email address ([a451f44](https://github.com/perfect-panel/ppanel-web/commit/a451f44))
- **forms**: Add step attribute to number inputs for better value control ([b8f4f1e](https://github.com/perfect-panel/ppanel-web/commit/b8f4f1e))
- **icon**: Comment out unused icon collection imports ([f17bf8d](https://github.com/perfect-panel/ppanel-web/commit/f17bf8d))
- **layout**: Remove unnecessary cookie initialization in Logout function ([3065c3a](https://github.com/perfect-panel/ppanel-web/commit/3065c3a))
- **locale**: Default value ([937408f](https://github.com/perfect-panel/ppanel-web/commit/937408f))
- **locale**: Document ([6f0fa20](https://github.com/perfect-panel/ppanel-web/commit/6f0fa20))
- **locale**: Empty ([3832d20](https://github.com/perfect-panel/ppanel-web/commit/3832d20))
- **locale**: Input Placeholder Webhook Domain ([bca0935](https://github.com/perfect-panel/ppanel-web/commit/bca0935))
- **locale**: Language Select ([0befdb0](https://github.com/perfect-panel/ppanel-web/commit/0befdb0))
- **locales**: Add error message for incorrect user information ([52c1d1f](https://github.com/perfect-panel/ppanel-web/commit/52c1d1f))
- **locales**: Add error message for incorrect user information ([3d92902](https://github.com/perfect-panel/ppanel-web/commit/3d92902))
- **locales**: Add logout message to authentication localization files ([1d0d911](https://github.com/perfect-panel/ppanel-web/commit/1d0d911))
- **locales**: Fixed description in multilingual files, updated text related to email registration functionality ([c356bc2](https://github.com/perfect-panel/ppanel-web/commit/c356bc2))
- **locales**: Order recharge related fields ([35210fe](https://github.com/perfect-panel/ppanel-web/commit/35210fe))
- **locales**: Removed language file import to clean up unnecessary language support ([68f6ab2](https://github.com/perfect-panel/ppanel-web/commit/68f6ab2))
- **locales**: Removed multilingual files to clean up unnecessary language support ([5b151cd](https://github.com/perfect-panel/ppanel-web/commit/5b151cd))
- **locale**: Subscription Path Description ([4c67387](https://github.com/perfect-panel/ppanel-web/commit/4c67387))
- **locales**: Update custom HTML description for clarity across multiple languages ([557c5cd](https://github.com/perfect-panel/ppanel-web/commit/557c5cd))
- **locales**: Update custom HTML description in language file, ([87381da](https://github.com/perfect-panel/ppanel-web/commit/87381da))
- **locales**: Update expiration time description from minutes to seconds in multiple languages ([5bac933](https://github.com/perfect-panel/ppanel-web/commit/5bac933))
- **locales**: Update Hong Kong ([6d0d069](https://github.com/perfect-panel/ppanel-web/commit/6d0d069))
- **locales**: Update invite code text to indicate it's optional ([6a34bfb](https://github.com/perfect-panel/ppanel-web/commit/6a34bfb))
- **logs**: Update log display to render key-value pairs and remove badge ([5ea6489](https://github.com/perfect-panel/ppanel-web/commit/5ea6489))
- **metadata**: Global metadata ([15d5ecf](https://github.com/perfect-panel/ppanel-web/commit/15d5ecf))
- **nav**: Comment out unused social login options to simplify navigation configuration ([cefcb31](https://github.com/perfect-panel/ppanel-web/commit/cefcb31))
- **node-config**: Add null checks for time slots and ensure proper handling of undefined values ([1cdb7e7](https://github.com/perfect-panel/ppanel-web/commit/1cdb7e7))
- **node**: Add country and city fields to the form schema and localization files ([8775fb6](https://github.com/perfect-panel/ppanel-web/commit/8775fb6))
- **node**: Handle potential null value for online users count ([fa2fb28](https://github.com/perfect-panel/ppanel-web/commit/fa2fb28))
- **node**: Locale and form ([38be4d5](https://github.com/perfect-panel/ppanel-web/commit/38be4d5))
- **node**: Port config ([a20834a](https://github.com/perfect-panel/ppanel-web/commit/a20834a))
- **node**: Reality config ([fadd17f](https://github.com/perfect-panel/ppanel-web/commit/fadd17f))
- **node**: Service Name config ([d0be685](https://github.com/perfect-panel/ppanel-web/commit/d0be685))
- **node**: TLS config ([57fae12](https://github.com/perfect-panel/ppanel-web/commit/57fae12))
- **node**: Trojan protocol config ([7e1eb90](https://github.com/perfect-panel/ppanel-web/commit/7e1eb90))
- **notify**: Ensure user info is updated after notification settings submission ([9bc3a94](https://github.com/perfect-panel/ppanel-web/commit/9bc3a94))
- **notify**: Set default values for notification settings to false ([3652819](https://github.com/perfect-panel/ppanel-web/commit/3652819))
- **oauth**: Refactor OAuth configuration types and update related API methods ([6227ba9](https://github.com/perfect-panel/ppanel-web/commit/6227ba9))
- **oauth**: Remove redundant checks when updating configuration to simplify logic ([9140b8a](https://github.com/perfect-panel/ppanel-web/commit/9140b8a))
- **payment**: Add notification URL field to payment management interface ([5c710e1](https://github.com/perfect-panel/ppanel-web/commit/5c710e1))
- **payment**: Config and types ([b0c87fb](https://github.com/perfect-panel/ppanel-web/commit/b0c87fb))
- **payment**: Fix payment related type definitions and update payment method references ([c3138a8](https://github.com/perfect-panel/ppanel-web/commit/c3138a8))
- **payment**: Qrcode ([a9a535b](https://github.com/perfect-panel/ppanel-web/commit/a9a535b))
- **payment**: Refactor payment form placeholder and update localization files ([4a4d364](https://github.com/perfect-panel/ppanel-web/commit/4a4d364))
- **payment**: Refactor purchaseCheckout usage and remove redundant code ([a5e2079](https://github.com/perfect-panel/ppanel-web/commit/a5e2079))
- **payment**: Replace window.open with window.location.href for checkout links ([1d8c765](https://github.com/perfect-panel/ppanel-web/commit/1d8c765))
- **payment**: Update checkout type from 'link' to 'url' for consistency ([136a1ab](https://github.com/perfect-panel/ppanel-web/commit/136a1ab))
- **payment**: Update payment information ([70d6a38](https://github.com/perfect-panel/ppanel-web/commit/70d6a38))
- **payment**: Update payment method update logic to include row data ([6752420](https://github.com/perfect-panel/ppanel-web/commit/6752420))
- **phone**: Update SMS expiration time field to use 'sms_expire_time' with default value of 300 ([18b07c7](https://github.com/perfect-panel/ppanel-web/commit/18b07c7))
- **profile**: Restore filter to ensure only valid OAuth accounts are shown ([315c8f9](https://github.com/perfect-panel/ppanel-web/commit/315c8f9))
- **purchasing**: Update payment type to lowercase and add optional chaining for discounts ([c06ea49](https://github.com/perfect-panel/ppanel-web/commit/c06ea49))
- **redirect**: Simplify redirect URL logic by removing unnecessary condition for sessionStorage ([c53ac61](https://github.com/perfect-panel/ppanel-web/commit/c53ac61))
- **redirect**: Update redirect URL logic to ensure proper handling of OAuth and auth paths ([7954762](https://github.com/perfect-panel/ppanel-web/commit/7954762))
- **register**: Adjust user email verification logic to handle domain suffix checks correctly ([686aa2d](https://github.com/perfect-panel/ppanel-web/commit/686aa2d))
- **request**: Add error code 40005 to trigger logout ([71bf002](https://github.com/perfect-panel/ppanel-web/commit/71bf002))
- **request**: Locale ([37d408f](https://github.com/perfect-panel/ppanel-web/commit/37d408f))
- **rule-form**: Remove redundant rule set display ([6e0c9b6](https://github.com/perfect-panel/ppanel-web/commit/6e0c9b6))
- **rules**: Remove unused MATCH rule ([674a01c](https://github.com/perfect-panel/ppanel-web/commit/674a01c))
- **site**: Add image upload functionality for site logo configuration ([4ea6e4a](https://github.com/perfect-panel/ppanel-web/commit/4ea6e4a))
- **site**: Se ref to store site configuration for updates ([0c8f091](https://github.com/perfect-panel/ppanel-web/commit/0c8f091))
- **sort**: Refactor sorting logic in NodeTable and SubscribeTable components for improved clarity and performance ([331bbea](https://github.com/perfect-panel/ppanel-web/commit/331bbea))
- **subscribe**: Add value prop to field in subscription form for proper state management ([328838d](https://github.com/perfect-panel/ppanel-web/commit/328838d))
- **subscribe**: Discount ([35a9f69](https://github.com/perfect-panel/ppanel-web/commit/35a9f69))
- **subscribe**: Extract Domain ([40d61a9](https://github.com/perfect-panel/ppanel-web/commit/40d61a9))
- **subscribe**: Handle optional values in price and discount calculations ([5939763](https://github.com/perfect-panel/ppanel-web/commit/5939763))
- **subscribe**: Jumps and internationalization ([13fdec3](https://github.com/perfect-panel/ppanel-web/commit/13fdec3))
- **subscribe**: Refactor discount calculations and default selection logic in subscription forms ([423b240](https://github.com/perfect-panel/ppanel-web/commit/423b240))
- **subscribe**: Server group id ([90e6764](https://github.com/perfect-panel/ppanel-web/commit/90e6764))
- **subscribe**: Update default selection logic in subscription form to ensure proper state management ([ef15374](https://github.com/perfect-panel/ppanel-web/commit/ef15374))
- **subscribe**: Update forms to include refetch functionality and improve toast messages ([fc55e95](https://github.com/perfect-panel/ppanel-web/commit/fc55e95))
- **subscribe**: Update payment return URL ([2b80496](https://github.com/perfect-panel/ppanel-web/commit/2b80496))
- **subscribe**: Update subscription domain placeholder to include examples; improve site name retrieval in global store ([c65a44c](https://github.com/perfect-panel/ppanel-web/commit/c65a44c))
- **subscribe**: Update value validation to check for number type in subscribe form ([6de29d5](https://github.com/perfect-panel/ppanel-web/commit/6de29d5))
- **subscription**: Add reset functionality for user subscription token ([39e89bf](https://github.com/perfect-panel/ppanel-web/commit/39e89bf))
- **table**: Update privacy policy tab translation key and remove unnecessary requestType from OAuth callback ([14b3af5](https://github.com/perfect-panel/ppanel-web/commit/14b3af5))
- **third-party-accounts**: Remove mobile display logic from third-party accounts component ([b4946f7](https://github.com/perfect-panel/ppanel-web/commit/b4946f7))
- **third-party-accounts**: Update redirect property name in binding response handling ([012e83a](https://github.com/perfect-panel/ppanel-web/commit/012e83a))
- **turnstile**: Turnstile_site_key ([0327b73](https://github.com/perfect-panel/ppanel-web/commit/0327b73))
- **type**: Fix ts type check error ([3cb0629](https://github.com/perfect-panel/ppanel-web/commit/3cb0629))
- **types**: Add 'gift_amount' field to API type definitions ([8f8a12a](https://github.com/perfect-panel/ppanel-web/commit/8f8a12a))
- **types**: Checking ([2992824](https://github.com/perfect-panel/ppanel-web/commit/2992824))
- **types**: Order type ([c7e50a9](https://github.com/perfect-panel/ppanel-web/commit/c7e50a9))
- **ui**: Bugs ([b023d0f](https://github.com/perfect-panel/ppanel-web/commit/b023d0f))
- **ui**: Components ([a7927d7](https://github.com/perfect-panel/ppanel-web/commit/a7927d7))
- **ui**: Fix json formatting ([e1ddd94](https://github.com/perfect-panel/ppanel-web/commit/e1ddd94))
- **ui**: Improve dashboard layout and enhance button functionality; open checkout URLs in a new tab ([fc0da76](https://github.com/perfect-panel/ppanel-web/commit/fc0da76))
- **ui**: Multiple display bugs ([f5d8fd3](https://github.com/perfect-panel/ppanel-web/commit/f5d8fd3))
- **user-nav**: Update user avatar and label to display telephone if email is not available ([7b6bb7b](https://github.com/perfect-panel/ppanel-web/commit/7b6bb7b))
- **user**: Add the 'gift_amount' field to the user service's type definition ([6301409](https://github.com/perfect-panel/ppanel-web/commit/6301409))
- **user**: Refactor user form validation and reset password fields ([6733fc2](https://github.com/perfect-panel/ppanel-web/commit/6733fc2))
- **user**: Update locales ([4e7d249](https://github.com/perfect-panel/ppanel-web/commit/4e7d249))
- **user**: Update notification and verify code settings ([574b043](https://github.com/perfect-panel/ppanel-web/commit/574b043))
- **user**: Update user identifier field and localizations ([1b6befa](https://github.com/perfect-panel/ppanel-web/commit/1b6befa))
- **user**: Update user subscribe display ([3bb714d](https://github.com/perfect-panel/ppanel-web/commit/3bb714d))
- **utils**: Login redirect url ([cbe5f0d](https://github.com/perfect-panel/ppanel-web/commit/cbe5f0d))
- More bugs ([2d88a3a](https://github.com/perfect-panel/ppanel-web/commit/2d88a3a))

### 👷 Build System

- **config**: Update pm2 config ([d95b425](https://github.com/perfect-panel/ppanel-web/commit/d95b425))

### 💄 Styles

- **dashboard**: Adjust grid layout and update image dimensions in application display ([f3204b7](https://github.com/perfect-panel/ppanel-web/commit/f3204b7))
- **dashboard**: Enhance card components with full height and improved empty state handling ([7e1d551](https://github.com/perfect-panel/ppanel-web/commit/7e1d551))
- **document**: Update ([0a8109b](https://github.com/perfect-panel/ppanel-web/commit/0a8109b))
- **globals**: Refactor delete confirmation button and update badge styles in node and subscribe tables ([30ae781](https://github.com/perfect-panel/ppanel-web/commit/30ae781))
- **locales**: Remove unused subscription labels from multiple locale files ([fb0c510](https://github.com/perfect-panel/ppanel-web/commit/fb0c510))
- **locales**: Update server.json to reorganize relay mode options and improve labels ([701cdee](https://github.com/perfect-panel/ppanel-web/commit/701cdee))
- **node**: Form ([d5f5add](https://github.com/perfect-panel/ppanel-web/commit/d5f5add))
- **node**: Improve layout and spacing in NodeStatusCell component ([136287d](https://github.com/perfect-panel/ppanel-web/commit/136287d))
- **node**: Protocol Tab ([2bcb925](https://github.com/perfect-panel/ppanel-web/commit/2bcb925))
- **time-slot**: Add chart display ([c44ad47](https://github.com/perfect-panel/ppanel-web/commit/c44ad47))
- **ui**: Update mobile style ([eda18bc](https://github.com/perfect-panel/ppanel-web/commit/eda18bc))
- Update node secret UI and add telephone code field to authentication form ([770932e](https://github.com/perfect-panel/ppanel-web/commit/770932e))

### 📝 Documentation

- **readme**: License name ([74cb16b](https://github.com/perfect-panel/ppanel-web/commit/74cb16b))

### 🔧 Continuous Integration

- **github**: Release docker ([5af60aa](https://github.com/perfect-panel/ppanel-web/commit/5af60aa))
- **step**: Update step name ([9eca618](https://github.com/perfect-panel/ppanel-web/commit/9eca618))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.34](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.33...v1.0.0-beta.34) (2025-04-02)

### ✨ Features

- **admin**: Add application and rule management entries to localization files ([8b43e69](https://github.com/perfect-panel/ppanel-web/commit/8b43e69))
- **api**: Add an interface to obtain user subscription details, update related type definitions and localized text ([cf5c39c](https://github.com/perfect-panel/ppanel-web/commit/cf5c39c))
- **user**: Integrate subscription list into user management, update request parameters and types ([8d49dac](https://github.com/perfect-panel/ppanel-web/commit/8d49dac))

### 🐛 Bug Fixes

- **admin**: Hidden versions and system upgrades ([64cd842](https://github.com/perfect-panel/ppanel-web/commit/64cd842))
- **admin**: Modify the label type in the rule form to a string array ([a7aa5fe](https://github.com/perfect-panel/ppanel-web/commit/a7aa5fe))
- **node**: Handle potential null value for online users count ([fa2fb28](https://github.com/perfect-panel/ppanel-web/commit/fa2fb28))
- **subscribe**: Add value prop to field in subscription form for proper state management ([328838d](https://github.com/perfect-panel/ppanel-web/commit/328838d))
- **subscribe**: Refactor discount calculations and default selection logic in subscription forms ([423b240](https://github.com/perfect-panel/ppanel-web/commit/423b240))
- **subscribe**: Update default selection logic in subscription form to ensure proper state management ([ef15374](https://github.com/perfect-panel/ppanel-web/commit/ef15374))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.33](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.32...v1.0.0-beta.33) (2025-03-18)

### 🐛 Bug Fixes

- **subscribe**: Handle optional values in price and discount calculations ([5939763](https://github.com/perfect-panel/ppanel-web/commit/5939763))

# [1.0.0-beta.32](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.31...v1.0.0-beta.32) (2025-03-17)

### 🐛 Bug Fixes

- **forms**: Add step attribute to number inputs for better value control ([b8f4f1e](https://github.com/perfect-panel/ppanel-web/commit/b8f4f1e))
- **locales**: Update invite code text to indicate it's optional ([6a34bfb](https://github.com/perfect-panel/ppanel-web/commit/6a34bfb))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.31](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.30...v1.0.0-beta.31) (2025-03-15)

### 🐛 Bug Fixes

- **site**: Se ref to store site configuration for updates ([0c8f091](https://github.com/perfect-panel/ppanel-web/commit/0c8f091))

# [1.0.0-beta.30](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.29...v1.0.0-beta.30) (2025-03-15)

### ✨ Features

- **api**: Add CheckoutOrder request and response types, and update user purchase request parameters ([48a1b97](https://github.com/perfect-panel/ppanel-web/commit/48a1b97))
- **email**: Add traffic exhaustion template ([bb3bd7b](https://github.com/perfect-panel/ppanel-web/commit/bb3bd7b))
- **formatting**: Update differenceInDays function to return whole days or two decimal places ([bf58f25](https://github.com/perfect-panel/ppanel-web/commit/bf58f25))
- **global**: Add custom data ([6dbebd1](https://github.com/perfect-panel/ppanel-web/commit/6dbebd1))
- **input**: Add minimum value constraint and enhance number handling in EnhancedInput ([ce31972](https://github.com/perfect-panel/ppanel-web/commit/ce31972))
- **loading**: Replace loading animation with a simpler spinner and loading text ([f72df3a](https://github.com/perfect-panel/ppanel-web/commit/f72df3a))
- **node-form**: Update number input fields to enforce step, min, and max values ([3f7b6d1](https://github.com/perfect-panel/ppanel-web/commit/3f7b6d1))
- **payment**: Add isEdit prop to PaymentForm and disable fields when editing ([85f55de](https://github.com/perfect-panel/ppanel-web/commit/85f55de))
- **timeline**: Simplify timeline component layout and remove commented-out code ([fbad3b0](https://github.com/perfect-panel/ppanel-web/commit/fbad3b0))

### 🎫 Chores

- **release**: V1.0.0-beta.27 [skip ci] ([092477b](https://github.com/perfect-panel/ppanel-web/commit/092477b))
- **release**: V1.0.0-beta.28 [skip ci] ([786ba0e](https://github.com/perfect-panel/ppanel-web/commit/786ba0e))
- Merge branch 'beta' into develop ([f219c52](https://github.com/perfect-panel/ppanel-web/commit/f219c52))

### 🐛 Bug Fixes

- **dashboard**: Update date display to use start_time if available ([e551232](https://github.com/perfect-panel/ppanel-web/commit/e551232))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.29](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.28...v1.0.0-beta.29) (2025-03-14)

### ✨ Features

- **api**: Add CheckoutOrder request and response types, and update user purchase request parameters ([dddc21c](https://github.com/perfect-panel/ppanel-web/commit/dddc21c))
- **loading**: Replace loading animation with a simpler spinner and loading text ([b8316bb](https://github.com/perfect-panel/ppanel-web/commit/b8316bb))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.28](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.27...v1.0.0-beta.28) (2025-03-13)

### ✨ Features

- **input**: Add minimum value constraint and enhance number handling in EnhancedInput ([94822d9](https://github.com/perfect-panel/ppanel-web/commit/94822d9))

# [1.0.0-beta.27](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.26...v1.0.0-beta.27) (2025-03-13)

### ♻ Code Refactoring

- **payment**: Reconstruct the payment page ([7109472](https://github.com/perfect-panel/ppanel-web/commit/7109472))
- Enhance user navigation dropdown ui and styling ([d2732e6](https://github.com/perfect-panel/ppanel-web/commit/d2732e6))

### ✨ Features

- **cdn**: Add CDN URL configuration and update related references ([0c90733](https://github.com/perfect-panel/ppanel-web/commit/0c90733))
- **payment**: Add bank card payment ([7fa3a57](https://github.com/perfect-panel/ppanel-web/commit/7fa3a57))
- **subscription**: Improve layout and organization of subscription detail tabs ([e4630f8](https://github.com/perfect-panel/ppanel-web/commit/e4630f8))
- **subscription**: Refactor subscription handling and update imports for better organization ([2215c7f](https://github.com/perfect-panel/ppanel-web/commit/2215c7f))

### 🎫 Chores

- **merge**: Bump version to 1.0.0-beta.26 and update changelog ([3222016](https://github.com/perfect-panel/ppanel-web/commit/3222016))

### 🐛 Bug Fixes

- **affiliate**: Update user identifier ([35f92c9](https://github.com/perfect-panel/ppanel-web/commit/35f92c9))
- **changelog**: Update change log style ([cfa3fc0](https://github.com/perfect-panel/ppanel-web/commit/cfa3fc0))
- **payment**: Add notification URL field to payment management interface ([5c710e1](https://github.com/perfect-panel/ppanel-web/commit/5c710e1))
- **payment**: Fix payment related type definitions and update payment method references ([c3138a8](https://github.com/perfect-panel/ppanel-web/commit/c3138a8))
- **payment**: Refactor purchaseCheckout usage and remove redundant code ([a5e2079](https://github.com/perfect-panel/ppanel-web/commit/a5e2079))
- **payment**: Update checkout type from 'link' to 'url' for consistency ([136a1ab](https://github.com/perfect-panel/ppanel-web/commit/136a1ab))
- **payment**: Update payment information ([70d6a38](https://github.com/perfect-panel/ppanel-web/commit/70d6a38))
- **payment**: Update payment method update logic to include row data ([6752420](https://github.com/perfect-panel/ppanel-web/commit/6752420))
- **purchasing**: Update payment type to lowercase and add optional chaining for discounts ([c06ea49](https://github.com/perfect-panel/ppanel-web/commit/c06ea49))
- **ui**: Improve dashboard layout and enhance button functionality; open checkout URLs in a new tab ([fc0da76](https://github.com/perfect-panel/ppanel-web/commit/fc0da76))
- **ui**: Multiple display bugs ([f5d8fd3](https://github.com/perfect-panel/ppanel-web/commit/f5d8fd3))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.28](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.27...v1.0.0-beta.28) (2025-03-13)

### ✨ Features

- **input**: Add minimum value constraint and enhance number handling in EnhancedInput ([94822d9](https://github.com/perfect-panel/ppanel-web/commit/94822d9))

# [1.0.0-beta.27](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.26...v1.0.0-beta.27) (2025-03-13)

### ♻ Code Refactoring

- **payment**: Reconstruct the payment page ([7109472](https://github.com/perfect-panel/ppanel-web/commit/7109472))
- Enhance user navigation dropdown ui and styling ([d2732e6](https://github.com/perfect-panel/ppanel-web/commit/d2732e6))

### ✨ Features

- **cdn**: Add CDN URL configuration and update related references ([0c90733](https://github.com/perfect-panel/ppanel-web/commit/0c90733))
- **payment**: Add bank card payment ([7fa3a57](https://github.com/perfect-panel/ppanel-web/commit/7fa3a57))
- **subscription**: Improve layout and organization of subscription detail tabs ([e4630f8](https://github.com/perfect-panel/ppanel-web/commit/e4630f8))
- **subscription**: Refactor subscription handling and update imports for better organization ([2215c7f](https://github.com/perfect-panel/ppanel-web/commit/2215c7f))

### 🎫 Chores

- **merge**: Bump version to 1.0.0-beta.26 and update changelog ([3222016](https://github.com/perfect-panel/ppanel-web/commit/3222016))

### 🐛 Bug Fixes

- **affiliate**: Update user identifier ([35f92c9](https://github.com/perfect-panel/ppanel-web/commit/35f92c9))
- **changelog**: Update change log style ([cfa3fc0](https://github.com/perfect-panel/ppanel-web/commit/cfa3fc0))
- **payment**: Add notification URL field to payment management interface ([5c710e1](https://github.com/perfect-panel/ppanel-web/commit/5c710e1))
- **payment**: Fix payment related type definitions and update payment method references ([c3138a8](https://github.com/perfect-panel/ppanel-web/commit/c3138a8))
- **payment**: Refactor purchaseCheckout usage and remove redundant code ([a5e2079](https://github.com/perfect-panel/ppanel-web/commit/a5e2079))
- **payment**: Update checkout type from 'link' to 'url' for consistency ([136a1ab](https://github.com/perfect-panel/ppanel-web/commit/136a1ab))
- **payment**: Update payment information ([70d6a38](https://github.com/perfect-panel/ppanel-web/commit/70d6a38))
- **payment**: Update payment method update logic to include row data ([6752420](https://github.com/perfect-panel/ppanel-web/commit/6752420))
- **purchasing**: Update payment type to lowercase and add optional chaining for discounts ([c06ea49](https://github.com/perfect-panel/ppanel-web/commit/c06ea49))
- **ui**: Improve dashboard layout and enhance button functionality; open checkout URLs in a new tab ([fc0da76](https://github.com/perfect-panel/ppanel-web/commit/fc0da76))
- **ui**: Multiple display bugs ([f5d8fd3](https://github.com/perfect-panel/ppanel-web/commit/f5d8fd3))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.26](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.25...v1.0.0-beta.26) (2025-03-02)

### 🐛 Bug Fixes

- **icon**: Comment out unused icon collection imports ([f17bf8d](https://github.com/perfect-panel/ppanel-web/commit/f17bf8d))

# [1.0.0-beta.25](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.24...v1.0.0-beta.25) (2025-03-01)

### ✨ Features

- **auth**: Add privacy policy link to the footer ([8e16ef1](https://github.com/perfect-panel/ppanel-web/commit/8e16ef1))

### 🐛 Bug Fixes

- **dashboard**: Display subscription creation date in user dashboard ([d0e6df0](https://github.com/perfect-panel/ppanel-web/commit/d0e6df0))
- **request**: Add error code 40005 to trigger logout ([71bf002](https://github.com/perfect-panel/ppanel-web/commit/71bf002))
- **subscribe**: Update payment return URL ([2b80496](https://github.com/perfect-panel/ppanel-web/commit/2b80496))

# [1.0.0-beta.24](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.23...v1.0.0-beta.24) (2025-02-27)

### ♻ Code Refactoring

- **ui**: Optimize document display ([2ca2992](https://github.com/perfect-panel/ppanel-web/commit/2ca2992))
- Reduce code complexity and improve readability ([e11f18c](https://github.com/perfect-panel/ppanel-web/commit/e11f18c))

### ✨ Features

- **loading**: Add loading components and integrate them in Providers ([d5847fa](https://github.com/perfect-panel/ppanel-web/commit/d5847fa))

### 🎫 Chores

- **merge**: Add advertising module and device settings ([0130e02](https://github.com/perfect-panel/ppanel-web/commit/0130e02))

### 🐛 Bug Fixes

- **locales**: Order recharge related fields ([35210fe](https://github.com/perfect-panel/ppanel-web/commit/35210fe))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.23](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.22...v1.0.0-beta.23) (2025-02-24)

### 🐛 Bug Fixes

- **auth**: Update email verification logic to use domain suffix check ([62662bb](https://github.com/perfect-panel/ppanel-web/commit/62662bb))

# [1.0.0-beta.22](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.21...v1.0.0-beta.22) (2025-02-23)

### 🐛 Bug Fixes

- **locales**: Removed language file import to clean up unnecessary language support ([68f6ab2](https://github.com/perfect-panel/ppanel-web/commit/68f6ab2))

# [1.0.0-beta.21](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.20...v1.0.0-beta.21) (2025-02-23)

### ✨ Features

- **privacy-policy**: Add privacy policy related text and links ([baa68f0](https://github.com/perfect-panel/ppanel-web/commit/baa68f0))

### 🐛 Bug Fixes

- **locales**: Removed multilingual files to clean up unnecessary language support ([5b151cd](https://github.com/perfect-panel/ppanel-web/commit/5b151cd))
- **locales**: Update custom HTML description in language file, ([87381da](https://github.com/perfect-panel/ppanel-web/commit/87381da))
- **table**: Update privacy policy tab translation key and remove unnecessary requestType from OAuth callback ([14b3af5](https://github.com/perfect-panel/ppanel-web/commit/14b3af5))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.20](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.19...v1.0.0-beta.20) (2025-02-23)

### 🐛 Bug Fixes

- **auth**: Add error handling to form submission and reset Turnstile on failure ([715d011](https://github.com/perfect-panel/ppanel-web/commit/715d011))
- **subscribe**: Update subscription domain placeholder to include examples; improve site name retrieval in global store ([c65a44c](https://github.com/perfect-panel/ppanel-web/commit/c65a44c))

# [1.0.0-beta.19](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.18...v1.0.0-beta.19) (2025-02-23)

### ✨ Features

- **form**: Make version field optional and set default value; update site domain placeholder for clarity ([42ba9e8](https://github.com/perfect-panel/ppanel-web/commit/42ba9e8))

### 🐛 Bug Fixes

- **auth**: Refactor forms to use Turnstile ref for reset functionality ([320a7dc](https://github.com/perfect-panel/ppanel-web/commit/320a7dc))

# [1.0.0-beta.18](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.17...v1.0.0-beta.18) (2025-02-22)

### ✨ Features

- **platform**: Update platform naming and add keywords and custom HTML fields ([6384237](https://github.com/perfect-panel/ppanel-web/commit/6384237))
- **site**: Added localization support for custom HTML and keyword fields ([f9d7736](https://github.com/perfect-panel/ppanel-web/commit/f9d7736))

### 🐛 Bug Fixes

- **locales**: Update custom HTML description for clarity across multiple languages ([557c5cd](https://github.com/perfect-panel/ppanel-web/commit/557c5cd))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.17](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.16...v1.0.0-beta.17) (2025-02-21)

### 🐛 Bug Fixes

- **auth**: Refactor reset password form to simplify code input and update placeholder text ([23833b4](https://github.com/perfect-panel/ppanel-web/commit/23833b4))

# [1.0.0-beta.16](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.15...v1.0.0-beta.16) (2025-02-21)

### 🐛 Bug Fixes

- **auth**: Simplify email verification code input rendering ([6f7bc37](https://github.com/perfect-panel/ppanel-web/commit/6f7bc37))

# [1.0.0-beta.15](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.14...v1.0.0-beta.15) (2025-02-21)

### 🐛 Bug Fixes

- **profile**: Restore filter to ensure only valid OAuth accounts are shown ([315c8f9](https://github.com/perfect-panel/ppanel-web/commit/315c8f9))

# [1.0.0-beta.14](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.13...v1.0.0-beta.14) (2025-02-21)

### 🐛 Bug Fixes

- **locales**: Fixed description in multilingual files, updated text related to email registration functionality ([c356bc2](https://github.com/perfect-panel/ppanel-web/commit/c356bc2))

# [1.0.0-beta.13](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.12...v1.0.0-beta.13) (2025-02-21)

### 🐛 Bug Fixes

- **api**: Fix type error in API request and add return URL parameter ([ee286dd](https://github.com/perfect-panel/ppanel-web/commit/ee286dd))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.12](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.11...v1.0.0-beta.12) (2025-02-18)

### 🐛 Bug Fixes

- More bugs ([2d88a3a](https://github.com/perfect-panel/ppanel-web/commit/2d88a3a))

# [1.0.0-beta.11](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.10...v1.0.0-beta.11) (2025-02-16)

### 🐛 Bug Fixes

- **email**: Update platform configuration handling to use current ref for consistency ([c90175b](https://github.com/perfect-panel/ppanel-web/commit/c90175b))

# [1.0.0-beta.10](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.9...v1.0.0-beta.10) (2025-02-16)

### 🐛 Bug Fixes

- **subscribe**: Update forms to include refetch functionality and improve toast messages ([fc55e95](https://github.com/perfect-panel/ppanel-web/commit/fc55e95))

# [1.0.0-beta.9](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.8...v1.0.0-beta.9) (2025-02-15)

### 🐛 Bug Fixes

- **register**: Adjust user email verification logic to handle domain suffix checks correctly ([686aa2d](https://github.com/perfect-panel/ppanel-web/commit/686aa2d))

# [1.0.0-beta.8](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.7...v1.0.0-beta.8) (2025-02-15)

### ✨ Features

- **accounts**: Update third-party account binding and unbinding ([1841552](https://github.com/perfect-panel/ppanel-web/commit/1841552))
- **auth-control**: Adding phone number labels to mobile verification configurations in multiple languages ([046740f](https://github.com/perfect-panel/ppanel-web/commit/046740f))
- **auth-control**: Update general ([3883646](https://github.com/perfect-panel/ppanel-web/commit/3883646))
- **auth**: Add type parameter to SendCode and update related API typings ([4198871](https://github.com/perfect-panel/ppanel-web/commit/4198871))
- **auth**: Refactor mobile authentication config to support whitelist functionality ([c761ec7](https://github.com/perfect-panel/ppanel-web/commit/c761ec7))
- **device**: Modify IMEI to device identifier support ([e3f9ef6](https://github.com/perfect-panel/ppanel-web/commit/e3f9ef6))
- **imei**: Add IMEI related internationalization support and menu items ([13c3337](https://github.com/perfect-panel/ppanel-web/commit/13c3337))
- **locales**: Add kick offline confirmation and success messages in multiple languages ([5db5343](https://github.com/perfect-panel/ppanel-web/commit/5db5343))
- **locales**: Update 'sms' to 'mobile' in authentication methods across multiple languages ([fea2171](https://github.com/perfect-panel/ppanel-web/commit/fea2171))
- **log**: Add message log retrieval functionality and update related typings ([1c0ecae](https://github.com/perfect-panel/ppanel-web/commit/1c0ecae))
- **profile**: Update localization strings and enhance third-party account binding ([2d1effb](https://github.com/perfect-panel/ppanel-web/commit/2d1effb))
- **release**: Extend supported platforms for Docker images, closes [#9](https://github.com/perfect-panel/ppanel-web/issues/9) ([e3a31eb](https://github.com/perfect-panel/ppanel-web/commit/e3a31eb))
- **subscription**: Add delete user subscription functionality ([1fc3a10](https://github.com/perfect-panel/ppanel-web/commit/1fc3a10))
- **ui**: Update input components and enhance card minimum width for better layout ([8a02310](https://github.com/perfect-panel/ppanel-web/commit/8a02310))
- **user**: Add user Detail ([3a3d223](https://github.com/perfect-panel/ppanel-web/commit/3a3d223))
- **user**: Add User Detail ([fdaf11b](https://github.com/perfect-panel/ppanel-web/commit/fdaf11b))

### 🐛 Bug Fixes

- **auth-control**: Fix citation error for platform values ([c940f3c](https://github.com/perfect-panel/ppanel-web/commit/c940f3c))
- **auth-control**: Fix citation error for platform values ([28813d2](https://github.com/perfect-panel/ppanel-web/commit/28813d2))
- **auth-control**: Rename phone_variable to phone_number in mobile verification configuration ([e5455aa](https://github.com/perfect-panel/ppanel-web/commit/e5455aa))
- **auth**: Update authentication configuration and localization strings ([47f2c58](https://github.com/perfect-panel/ppanel-web/commit/47f2c58))
- **locales**: Update expiration time description from minutes to seconds in multiple languages ([5bac933](https://github.com/perfect-panel/ppanel-web/commit/5bac933))
- **notify**: Ensure user info is updated after notification settings submission ([9bc3a94](https://github.com/perfect-panel/ppanel-web/commit/9bc3a94))
- **notify**: Set default values for notification settings to false ([3652819](https://github.com/perfect-panel/ppanel-web/commit/3652819))
- **third-party-accounts**: Remove mobile display logic from third-party accounts component ([b4946f7](https://github.com/perfect-panel/ppanel-web/commit/b4946f7))
- **third-party-accounts**: Update redirect property name in binding response handling ([012e83a](https://github.com/perfect-panel/ppanel-web/commit/012e83a))
- **user**: Refactor user form validation and reset password fields ([6733fc2](https://github.com/perfect-panel/ppanel-web/commit/6733fc2))
- **user**: Update locales ([4e7d249](https://github.com/perfect-panel/ppanel-web/commit/4e7d249))
- **user**: Update notification and verify code settings ([574b043](https://github.com/perfect-panel/ppanel-web/commit/574b043))

### 💄 Styles

- **dashboard**: Enhance card components with full height and improved empty state handling ([7e1d551](https://github.com/perfect-panel/ppanel-web/commit/7e1d551))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.7](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.6...v1.0.0-beta.7) (2025-01-28)

### ♻ Code Refactoring

- **sbscribe**: Rename and reorganize components for better structure and clarity ([5e5e4ed](https://github.com/perfect-panel/ppanel-web/commit/5e5e4ed))

### ✨ Features

- **auth**: Add email and SMS code sending functionality with localization updates ([57eaa55](https://github.com/perfect-panel/ppanel-web/commit/57eaa55))
- **auth**: Add Oauth configuration for Telegram, Facebook, Google, Github, and Apple ([18ee600](https://github.com/perfect-panel/ppanel-web/commit/18ee600))
- **auth**: Add SMS and email configuration options to global store and update localization ([4acf7b1](https://github.com/perfect-panel/ppanel-web/commit/4acf7b1))
- **auth**: Enhance user registration with invite handling and logo display ([207bc24](https://github.com/perfect-panel/ppanel-web/commit/207bc24))
- **auth**: Redirect user after OAuth login and add logos icon collection ([aa6dda8](https://github.com/perfect-panel/ppanel-web/commit/aa6dda8))
- **config**: Add application selection and encryption settings to configuration form ([88b3504](https://github.com/perfect-panel/ppanel-web/commit/88b3504))
- **config**: Update encryption fields in configuration form and refactor OAuth callback parameters ([652e032](https://github.com/perfect-panel/ppanel-web/commit/652e032))
- **global**: Add SMS configuration options to global store ([39a9ce6](https://github.com/perfect-panel/ppanel-web/commit/39a9ce6))
- **locales**: Add area code and telephone fields to user forms in multiple languages ([9b8258c](https://github.com/perfect-panel/ppanel-web/commit/9b8258c))
- **locales**: Add description information of communication keys and encryption methods to enhance client configuration capabilities ([d1f5a9b](https://github.com/perfect-panel/ppanel-web/commit/d1f5a9b))
- **node**: Add tags ([f408fdf](https://github.com/perfect-panel/ppanel-web/commit/f408fdf))
- **node**: Move the node configuration to the server module ([7f0f5ce](https://github.com/perfect-panel/ppanel-web/commit/7f0f5ce))
- **oauth**: Add certification component for handling OAuth login callbacks and improve user authentication flow ([5ed04c0](https://github.com/perfect-panel/ppanel-web/commit/5ed04c0))
- **oauth**: Implement OAuth token retrieval and refactor login callback handling ([40a6f7c](https://github.com/perfect-panel/ppanel-web/commit/40a6f7c))
- **oauth**: Refactor platform parameter handling and improve logout redirection logic ([8346c85](https://github.com/perfect-panel/ppanel-web/commit/8346c85))
- **oauth**: Update OAuth login handling to use callback parameter and improve URL parameter retrieval ([9227411](https://github.com/perfect-panel/ppanel-web/commit/9227411))
- **sms**: Update locales ([938363b](https://github.com/perfect-panel/ppanel-web/commit/938363b))
- **subscribe**: Add 'sold' column to SubscribeTable and update inventory terminology ([19619fd](https://github.com/perfect-panel/ppanel-web/commit/19619fd))
- **subscribe**: Move subscription configuration and application to subscription module ([f90d4d2](https://github.com/perfect-panel/ppanel-web/commit/f90d4d2))
- **subscribe**: Update SubscribeTable component to use API.SubscribeItem type and ensure proper type casting ([f26f1c2](https://github.com/perfect-panel/ppanel-web/commit/f26f1c2))
- **tutorial**: Fetch the latest tutorial version from GitHub API for dynamic URL generation ([28f8c78](https://github.com/perfect-panel/ppanel-web/commit/28f8c78))
- **user**: Add telephone input with area code selection and update localization ([585b99c](https://github.com/perfect-panel/ppanel-web/commit/585b99c))
- Update Auth Control ([c59742a](https://github.com/perfect-panel/ppanel-web/commit/c59742a))

### 🐛 Bug Fixes

- **api**: Remove redundant requestType parameter in appleLoginCallback ([0aa5d5b](https://github.com/perfect-panel/ppanel-web/commit/0aa5d5b))
- **api**: Rename app-related functions and types to application for consistency ([9d8b814](https://github.com/perfect-panel/ppanel-web/commit/9d8b814))
- **api**: Update subscription_protocol to subscribe_type for consistency across services ([b6da51b](https://github.com/perfect-panel/ppanel-web/commit/b6da51b))
- **auth**: Change Textarea value to defaultValue for client_secret in Apple auth page ([69fc670](https://github.com/perfect-panel/ppanel-web/commit/69fc670))
- **auth**: Remove unused telephone code login function and update typings for telephone login requests ([7239685](https://github.com/perfect-panel/ppanel-web/commit/7239685))
- **auth**: Require minimum length for invite string when forced invite is enabled ([a604f28](https://github.com/perfect-panel/ppanel-web/commit/a604f28))
- **auth**: Update user authentication flow to include email and phone code verification ([5d078fd](https://github.com/perfect-panel/ppanel-web/commit/5d078fd))
- **dashboard**: Improve URL encoding for subscription links and enhance success message handling ([4983c33](https://github.com/perfect-panel/ppanel-web/commit/4983c33))
- **dashboard**: Update icon imports for platform consistency and adjust icon size ([3e8912e](https://github.com/perfect-panel/ppanel-web/commit/3e8912e))
- **dashboard**: Update platform detection logic and improve layout responsiveness ([b0aa364](https://github.com/perfect-panel/ppanel-web/commit/b0aa364))
- **deps**: Remove outdated @iconify/react dependency and add iconify-json packages ([d6fbc38](https://github.com/perfect-panel/ppanel-web/commit/d6fbc38))
- **locales**: Add error message for incorrect user information ([52c1d1f](https://github.com/perfect-panel/ppanel-web/commit/52c1d1f))
- **locales**: Add error message for incorrect user information ([3d92902](https://github.com/perfect-panel/ppanel-web/commit/3d92902))
- **locales**: Add logout message to authentication localization files ([1d0d911](https://github.com/perfect-panel/ppanel-web/commit/1d0d911))
- **nav**: Comment out unused social login options to simplify navigation configuration ([cefcb31](https://github.com/perfect-panel/ppanel-web/commit/cefcb31))
- **node-config**: Add null checks for time slots and ensure proper handling of undefined values ([1cdb7e7](https://github.com/perfect-panel/ppanel-web/commit/1cdb7e7))
- **node**: Add country and city fields to the form schema and localization files ([8775fb6](https://github.com/perfect-panel/ppanel-web/commit/8775fb6))
- **oauth**: Refactor OAuth configuration types and update related API methods ([6227ba9](https://github.com/perfect-panel/ppanel-web/commit/6227ba9))
- **oauth**: Remove redundant checks when updating configuration to simplify logic ([9140b8a](https://github.com/perfect-panel/ppanel-web/commit/9140b8a))
- **payment**: Replace window.open with window.location.href for checkout links ([1d8c765](https://github.com/perfect-panel/ppanel-web/commit/1d8c765))
- **phone**: Update SMS expiration time field to use 'sms_expire_time' with default value of 300 ([18b07c7](https://github.com/perfect-panel/ppanel-web/commit/18b07c7))
- **redirect**: Simplify redirect URL logic by removing unnecessary condition for sessionStorage ([c53ac61](https://github.com/perfect-panel/ppanel-web/commit/c53ac61))
- **redirect**: Update redirect URL logic to ensure proper handling of OAuth and auth paths ([7954762](https://github.com/perfect-panel/ppanel-web/commit/7954762))
- **site**: Add image upload functionality for site logo configuration ([4ea6e4a](https://github.com/perfect-panel/ppanel-web/commit/4ea6e4a))
- **sort**: Refactor sorting logic in NodeTable and SubscribeTable components for improved clarity and performance ([331bbea](https://github.com/perfect-panel/ppanel-web/commit/331bbea))
- **subscription**: Add reset functionality for user subscription token ([39e89bf](https://github.com/perfect-panel/ppanel-web/commit/39e89bf))
- **type**: Fix ts type check error ([3cb0629](https://github.com/perfect-panel/ppanel-web/commit/3cb0629))
- **user-nav**: Update user avatar and label to display telephone if email is not available ([7b6bb7b](https://github.com/perfect-panel/ppanel-web/commit/7b6bb7b))
- **user**: Update user identifier field and localizations ([1b6befa](https://github.com/perfect-panel/ppanel-web/commit/1b6befa))

### 💄 Styles

- **dashboard**: Adjust grid layout and update image dimensions in application display ([f3204b7](https://github.com/perfect-panel/ppanel-web/commit/f3204b7))
- **globals**: Refactor delete confirmation button and update badge styles in node and subscribe tables ([30ae781](https://github.com/perfect-panel/ppanel-web/commit/30ae781))
- Update node secret UI and add telephone code field to authentication form ([770932e](https://github.com/perfect-panel/ppanel-web/commit/770932e))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.6](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.5...v1.0.0-beta.6) (2025-01-10)

### 🐛 Bug Fixes

- **auth**: Update UserCheckForm to use setInitialValues and modify onSubmit type ([c984c0d](https://github.com/perfect-panel/ppanel-web/commit/c984c0d))

# [1.0.0-beta.5](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.4...v1.0.0-beta.5) (2025-01-09)

### ✨ Features

- **locales**: Replace 'nodeGroupId' with 'groupId' in multiple language files for consistency ([a4e9d5d](https://github.com/perfect-panel/ppanel-web/commit/a4e9d5d))
- **locales**: Update 'deductBalance' to 'giftAmount' across multiple languages and fix newline in announcement.json ([70497af](https://github.com/perfect-panel/ppanel-web/commit/70497af))
- **stats**: Replace dynamic stat fetching with environment constants for user, server, and location counts ([46ae166](https://github.com/perfect-panel/ppanel-web/commit/46ae166))
- **subscribe**: Update suffix from 'MB' to 'Mbps' and enhance speed limit display logic ([3547bb1](https://github.com/perfect-panel/ppanel-web/commit/3547bb1))
- **user**: Add 'gift_amount' field and update related references in user services and components ([b13c77e](https://github.com/perfect-panel/ppanel-web/commit/b13c77e))

### 🎫 Chores

- **deps**: Update package dependencies across multiple projects for improved stability and performance ([b01a5bc](https://github.com/perfect-panel/ppanel-web/commit/b01a5bc))
- **ui**: Update package dependencies for improved stability and performance ([25da429](https://github.com/perfect-panel/ppanel-web/commit/25da429))

### 🐛 Bug Fixes

- **api**: Replace 'deduction' with 'gift_amount' and add 'commission' field in type definitions ([77edf1d](https://github.com/perfect-panel/ppanel-web/commit/77edf1d))
- **api**: Update API type definitions to replace 'deduction' with 'gift_amount' and make 'commission' field optional ([c2af060](https://github.com/perfect-panel/ppanel-web/commit/c2af060))
- **auth**: Refactor user authentication forms to remove global store dependency and improve type handling ([12026b0](https://github.com/perfect-panel/ppanel-web/commit/12026b0))
- **coupon**: Rename 'server' field to 'subscribe' in coupon form and update coupon update request type ([f8b6d82](https://github.com/perfect-panel/ppanel-web/commit/f8b6d82))
- **types**: Add 'gift_amount' field to API type definitions ([8f8a12a](https://github.com/perfect-panel/ppanel-web/commit/8f8a12a))
- **user**: Add the 'gift_amount' field to the user service's type definition ([6301409](https://github.com/perfect-panel/ppanel-web/commit/6301409))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.4](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.3...v1.0.0-beta.4) (2025-01-07)

### ♻ Code Refactoring

- **auth**: Refactor user authorization handling and improve error logging ([68bc18f](https://github.com/perfect-panel/ppanel-web/commit/68bc18f))

### ✨ Features

- **affiliate**: Add Affiliate component with commission display and invite link functionality ([4aea4e8](https://github.com/perfect-panel/ppanel-web/commit/4aea4e8))
- **affiliate**: Update affiliate component to display total commission and improve data fetching ([cc834ca](https://github.com/perfect-panel/ppanel-web/commit/cc834ca))
- **api**: Add new subscription properties and locale support for deduction ratios and reset cycles ([fec80f5](https://github.com/perfect-panel/ppanel-web/commit/fec80f5))
- **api**: Add Time Period Configuration ([837157c](https://github.com/perfect-panel/ppanel-web/commit/837157c))
- **favicon**: Update SVG favicon design for admin and user interfaces ([1d91738](https://github.com/perfect-panel/ppanel-web/commit/1d91738))
- **node**: Add serverKey ([25ce37e](https://github.com/perfect-panel/ppanel-web/commit/25ce37e))
- **relay**: Add relay mode configuration and update related schemas ([3cc9477](https://github.com/perfect-panel/ppanel-web/commit/3cc9477))
- **schema**: Add security field to hysteria2 and tuic schemas ([cd59d44](https://github.com/perfect-panel/ppanel-web/commit/cd59d44))
- **subscribe**: Add reset_time to API typings and update unsubscribe logic ([eeea165](https://github.com/perfect-panel/ppanel-web/commit/eeea165))
- **subscribe**: Add subscribe_discount type ([f99c604](https://github.com/perfect-panel/ppanel-web/commit/f99c604))
- **subscribe**: Add subscription credits ([5bc7905](https://github.com/perfect-panel/ppanel-web/commit/5bc7905))
- **subscribe**: Add unsubscribe functionality with confirmation messages and localized strings ([b2a2f42](https://github.com/perfect-panel/ppanel-web/commit/b2a2f42))
- **subscribe**: Improve error handling in subscription forms and update component props ([d28a10b](https://github.com/perfect-panel/ppanel-web/commit/d28a10b))
- **subscribe**: Improve layout and styling in subscription components ([5766376](https://github.com/perfect-panel/ppanel-web/commit/5766376))
- **subscription**: Add localized messages for existing subscriptions and deletion restrictions ([e8a72d5](https://github.com/perfect-panel/ppanel-web/commit/e8a72d5))

### 🎫 Chores

- Update changelog, enhance prepare script, and add openapi command ([a93db4e](https://github.com/perfect-panel/ppanel-web/commit/a93db4e))

### 🐛 Bug Fixes

- **dashboard**: Correct progress value calculations and update groupId accessor ([36c7667](https://github.com/perfect-panel/ppanel-web/commit/36c7667))
- **layout**: Remove unnecessary cookie initialization in Logout function ([3065c3a](https://github.com/perfect-panel/ppanel-web/commit/3065c3a))
- **locales**: Update Hong Kong ([6d0d069](https://github.com/perfect-panel/ppanel-web/commit/6d0d069))
- **subscribe**: Update value validation to check for number type in subscribe form ([6de29d5](https://github.com/perfect-panel/ppanel-web/commit/6de29d5))

### 💄 Styles

- **locales**: Remove unused subscription labels from multiple locale files ([fb0c510](https://github.com/perfect-panel/ppanel-web/commit/fb0c510))
- **locales**: Update server.json to reorganize relay mode options and improve labels ([701cdee](https://github.com/perfect-panel/ppanel-web/commit/701cdee))
- **node**: Improve layout and spacing in NodeStatusCell component ([136287d](https://github.com/perfect-panel/ppanel-web/commit/136287d))
- **time-slot**: Add chart display ([c44ad47](https://github.com/perfect-panel/ppanel-web/commit/c44ad47))

<a name="readme-top"></a>

# Changelog

# [1.0.0-beta.3](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.2...v1.0.0-beta.3) (2024-12-26)

### ♻ Code Refactoring

- **publish**: Simplify environment variable handling and improve build script ([88ea21b](https://github.com/perfect-panel/ppanel-web/commit/88ea21b))

# [1.0.0-beta.2](https://github.com/perfect-panel/ppanel-web/compare/v1.0.0-beta.1...v1.0.0-beta.2) (2024-12-26)

### 👷 Build System

- **config**: Update pm2 config ([d95b425](https://github.com/perfect-panel/ppanel-web/commit/d95b425))

### 🔧 Continuous Integration

- **step**: Update step name ([9eca618](https://github.com/perfect-panel/ppanel-web/commit/9eca618))

# 1.0.0-beta.1 (2024-12-25)

### ♻ Code Refactoring

- **api**: Sort and Announcement ([38d5616](https://github.com/perfect-panel/ppanel-web/commit/38d5616))
- **config**: GenerateMetadata ([a0bb101](https://github.com/perfect-panel/ppanel-web/commit/a0bb101))
- **config**: Viewport ([24b8601](https://github.com/perfect-panel/ppanel-web/commit/24b8601))
- **core**: Restructure project for better module separation ([9d0cb8b](https://github.com/perfect-panel/ppanel-web/commit/9d0cb8b))
- **deps**: Update ([19837a1](https://github.com/perfect-panel/ppanel-web/commit/19837a1))
- **empty**: Content ([aa4c667](https://github.com/perfect-panel/ppanel-web/commit/aa4c667))
- **ui**: Dependencies ([727d779](https://github.com/perfect-panel/ppanel-web/commit/727d779))
- **ui**: Layout ([9262d7d](https://github.com/perfect-panel/ppanel-web/commit/9262d7d))

### ⚡ Performance Improvements

- **subscribe**: Form discount price ([059a892](https://github.com/perfect-panel/ppanel-web/commit/059a892))

### ✨ Features

- **ad**: Advertise ([b1105cd](https://github.com/perfect-panel/ppanel-web/commit/b1105cd))
- **affiliate**: Affiliate Detail ([a782c17](https://github.com/perfect-panel/ppanel-web/commit/a782c17))
- **affiliate**: Commission Rate ([5eec430](https://github.com/perfect-panel/ppanel-web/commit/5eec430))
- **announcement**: Popup and pinned ([f3680a7](https://github.com/perfect-panel/ppanel-web/commit/f3680a7))
- **api**: Telegram ([17ce96a](https://github.com/perfect-panel/ppanel-web/commit/17ce96a))
- **billing**: Update Billing ([078fc9d](https://github.com/perfect-panel/ppanel-web/commit/078fc9d))
- **config**: FormatBytes ([9251a09](https://github.com/perfect-panel/ppanel-web/commit/9251a09))
- **config**: Protocol type ([a3b45b4](https://github.com/perfect-panel/ppanel-web/commit/a3b45b4))
- **config**: Webhook Domain ([01e06c6](https://github.com/perfect-panel/ppanel-web/commit/01e06c6))
- **dashboard**: Optimization ([5b3f4b4](https://github.com/perfect-panel/ppanel-web/commit/5b3f4b4))
- **dashboard**: Statistics ([2926abc](https://github.com/perfect-panel/ppanel-web/commit/2926abc))
- **header**: Update locales ([bfb6c27](https://github.com/perfect-panel/ppanel-web/commit/bfb6c27))
- **locale**: Add Persian ([93a0a88](https://github.com/perfect-panel/ppanel-web/commit/93a0a88))
- **locales**: Add multiple languages ([b243ab9](https://github.com/perfect-panel/ppanel-web/commit/b243ab9))
- **node-subscription**: Add copy functionality for columns ([3a81e37](https://github.com/perfect-panel/ppanel-web/commit/3a81e37))
- **node**: Add NodeStatus ([c712624](https://github.com/perfect-panel/ppanel-web/commit/c712624))
- **node**: Add protocol ([301b635](https://github.com/perfect-panel/ppanel-web/commit/301b635))
- **node**: Add status ([c06372b](https://github.com/perfect-panel/ppanel-web/commit/c06372b))
- **subscribe**: Add unit time ([39d07ec](https://github.com/perfect-panel/ppanel-web/commit/39d07ec))
- **table**: Add sorting support for Node and subscription columns ([27924b0](https://github.com/perfect-panel/ppanel-web/commit/27924b0))
- **table**: Supports drag and drop sorting ([2f56ef5](https://github.com/perfect-panel/ppanel-web/commit/2f56ef5))
- **tos**: Display data ([6024454](https://github.com/perfect-panel/ppanel-web/commit/6024454))
- **tutorial**: Add common tutorial list ([872252c](https://github.com/perfect-panel/ppanel-web/commit/872252c))
- **ui**: System Tool ([1836980](https://github.com/perfect-panel/ppanel-web/commit/1836980))
- **ui**: Update homepage data ([8425b13](https://github.com/perfect-panel/ppanel-web/commit/8425b13))

### 🎫 Chores

- **config**: Entry locale ([5737331](https://github.com/perfect-panel/ppanel-web/commit/5737331))
- **init**: Project initialization ([829edfa](https://github.com/perfect-panel/ppanel-web/commit/829edfa))

### 🐛 Bug Fixes

- **api**: PreCreateOrder ([ca747f5](https://github.com/perfect-panel/ppanel-web/commit/ca747f5))
- **api**: Purge ([98c1c30](https://github.com/perfect-panel/ppanel-web/commit/98c1c30))
- **api**: Server and order ([255bd82](https://github.com/perfect-panel/ppanel-web/commit/255bd82))
- **api**: Statistics ([7962162](https://github.com/perfect-panel/ppanel-web/commit/7962162))
- **api**: Subscribe token ([1932ba7](https://github.com/perfect-panel/ppanel-web/commit/1932ba7))
- **api**: Update Model ([39aaa73](https://github.com/perfect-panel/ppanel-web/commit/39aaa73))
- **billing**: ExpiryDate ([e85e545](https://github.com/perfect-panel/ppanel-web/commit/e85e545))
- **billing**: I18n and styles ([81e0f21](https://github.com/perfect-panel/ppanel-web/commit/81e0f21))
- **config**: AlipayF2F ([6c07107](https://github.com/perfect-panel/ppanel-web/commit/6c07107))
- **config**: Bugs ([f57e40c](https://github.com/perfect-panel/ppanel-web/commit/f57e40c))
- **config**: Checkout Order ([a31e763](https://github.com/perfect-panel/ppanel-web/commit/a31e763))
- **config**: FormatBytes ([bbc2da0](https://github.com/perfect-panel/ppanel-web/commit/bbc2da0))
- **config**: NoStore ([2cc18cf](https://github.com/perfect-panel/ppanel-web/commit/2cc18cf))
- **config**: Runtime env ([a1e4999](https://github.com/perfect-panel/ppanel-web/commit/a1e4999))
- **config**: Status Percentag ([8f322fb](https://github.com/perfect-panel/ppanel-web/commit/8f322fb))
- **config**: SubLink ([1c61966](https://github.com/perfect-panel/ppanel-web/commit/1c61966))
- **config**: Subscribe Link ([11ea821](https://github.com/perfect-panel/ppanel-web/commit/11ea821))
- **controller**: Order status ([8c6a097](https://github.com/perfect-panel/ppanel-web/commit/8c6a097))
- **dashboard**: Format Bytes ([d8b0bd9](https://github.com/perfect-panel/ppanel-web/commit/d8b0bd9))
- **deps**: Typescript config ([34e24b8](https://github.com/perfect-panel/ppanel-web/commit/34e24b8))
- **deps**: Update clipboard ([5572710](https://github.com/perfect-panel/ppanel-web/commit/5572710))
- **editor**: Change value ([4fdfeb2](https://github.com/perfect-panel/ppanel-web/commit/4fdfeb2))
- **footer**: Email address ([a451f44](https://github.com/perfect-panel/ppanel-web/commit/a451f44))
- **locale**: Default value ([937408f](https://github.com/perfect-panel/ppanel-web/commit/937408f))
- **locale**: Document ([6f0fa20](https://github.com/perfect-panel/ppanel-web/commit/6f0fa20))
- **locale**: Empty ([3832d20](https://github.com/perfect-panel/ppanel-web/commit/3832d20))
- **locale**: Input Placeholder Webhook Domain ([bca0935](https://github.com/perfect-panel/ppanel-web/commit/bca0935))
- **locale**: Language Select ([0befdb0](https://github.com/perfect-panel/ppanel-web/commit/0befdb0))
- **locale**: Subscription Path Description ([4c67387](https://github.com/perfect-panel/ppanel-web/commit/4c67387))
- **metadata**: Global metadata ([15d5ecf](https://github.com/perfect-panel/ppanel-web/commit/15d5ecf))
- **node**: Locale and form ([38be4d5](https://github.com/perfect-panel/ppanel-web/commit/38be4d5))
- **node**: Port config ([a20834a](https://github.com/perfect-panel/ppanel-web/commit/a20834a))
- **node**: Reality config ([fadd17f](https://github.com/perfect-panel/ppanel-web/commit/fadd17f))
- **node**: Service Name config ([d0be685](https://github.com/perfect-panel/ppanel-web/commit/d0be685))
- **node**: TLS config ([57fae12](https://github.com/perfect-panel/ppanel-web/commit/57fae12))
- **node**: Trojan protocol config ([7e1eb90](https://github.com/perfect-panel/ppanel-web/commit/7e1eb90))
- **payment**: Config and types ([b0c87fb](https://github.com/perfect-panel/ppanel-web/commit/b0c87fb))
- **payment**: Qrcode ([a9a535b](https://github.com/perfect-panel/ppanel-web/commit/a9a535b))
- **request**: Locale ([37d408f](https://github.com/perfect-panel/ppanel-web/commit/37d408f))
- **subscribe**: Discount ([35a9f69](https://github.com/perfect-panel/ppanel-web/commit/35a9f69))
- **subscribe**: Extract Domain ([40d61a9](https://github.com/perfect-panel/ppanel-web/commit/40d61a9))
- **subscribe**: Jumps and internationalization ([13fdec3](https://github.com/perfect-panel/ppanel-web/commit/13fdec3))
- **subscribe**: Server group id ([90e6764](https://github.com/perfect-panel/ppanel-web/commit/90e6764))
- **turnstile**: Turnstile_site_key ([0327b73](https://github.com/perfect-panel/ppanel-web/commit/0327b73))
- **types**: Checking ([2992824](https://github.com/perfect-panel/ppanel-web/commit/2992824))
- **types**: Order type ([c7e50a9](https://github.com/perfect-panel/ppanel-web/commit/c7e50a9))
- **ui**: Bugs ([b023d0f](https://github.com/perfect-panel/ppanel-web/commit/b023d0f))
- **ui**: Components ([a7927d7](https://github.com/perfect-panel/ppanel-web/commit/a7927d7))
- **ui**: Fix json formatting ([e1ddd94](https://github.com/perfect-panel/ppanel-web/commit/e1ddd94))
- **utils**: Login redirect url ([cbe5f0d](https://github.com/perfect-panel/ppanel-web/commit/cbe5f0d))

### 💄 Styles

- **document**: Update ([0a8109b](https://github.com/perfect-panel/ppanel-web/commit/0a8109b))
- **node**: Form ([d5f5add](https://github.com/perfect-panel/ppanel-web/commit/d5f5add))
- **node**: Protocol Tab ([2bcb925](https://github.com/perfect-panel/ppanel-web/commit/2bcb925))
- **ui**: Update mobile style ([eda18bc](https://github.com/perfect-panel/ppanel-web/commit/eda18bc))

### 📝 Documentation

- **readme**: License name ([74cb16b](https://github.com/perfect-panel/ppanel-web/commit/74cb16b))

### 🔧 Continuous Integration

- **github**: Release docker ([5af60aa](https://github.com/perfect-panel/ppanel-web/commit/5af60aa))
