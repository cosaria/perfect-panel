/**
 * A shared CommitLint configuration for the repository.
 *
 * @type {import("commitlint").Config}
 * */
module.exports = {
  $schema: 'https://json.schemastore.org/commitlintrc',
  extends: ['gitmoji'],
  rules: {
    'footer-leading-blank': [0, 'never'],
    'header-max-length': [0, 'never'],
  },
};
