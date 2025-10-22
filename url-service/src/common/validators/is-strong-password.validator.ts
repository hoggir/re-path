import {
  registerDecorator,
  ValidationOptions,
  ValidationArguments,
} from 'class-validator';

export interface StrongPasswordOptions {
  minLength?: number;
  requireUppercase?: boolean;
  requireLowercase?: boolean;
  requireNumbers?: boolean;
  requireSpecialChars?: boolean;
  allowWhitespace?: boolean;
}

export function IsStrongPassword(
  options?: StrongPasswordOptions,
  validationOptions?: ValidationOptions,
) {
  return function (object: object, propertyName: string) {
    registerDecorator({
      name: 'isStrongPassword',
      target: object.constructor,
      propertyName: propertyName,
      options: validationOptions,
      constraints: [options],
      validator: {
        validate(value: any, args: ValidationArguments) {
          if (typeof value !== 'string') {
            return false;
          }

          const opts: StrongPasswordOptions = args.constraints[0] || {};
          const minLength = opts.minLength ?? 8;
          const requireUppercase = opts.requireUppercase ?? true;
          const requireLowercase = opts.requireLowercase ?? true;
          const requireNumbers = opts.requireNumbers ?? true;
          const requireSpecialChars = opts.requireSpecialChars ?? true;
          const allowWhitespace = opts.allowWhitespace ?? false;

          if (value.length < minLength) {
            return false;
          }

          if (!allowWhitespace && /\s/.test(value)) {
            return false;
          }

          if (requireUppercase && !/[A-Z]/.test(value)) {
            return false;
          }

          if (requireLowercase && !/[a-z]/.test(value)) {
            return false;
          }

          if (requireNumbers && !/\d/.test(value)) {
            return false;
          }

          if (
            requireSpecialChars &&
            !/[@$!%*?&#^()_+\-=\[\]{};':"\\|,.<>\/]/.test(value)
          ) {
            return false;
          }

          return true;
        },
        defaultMessage(args: ValidationArguments) {
          const opts: StrongPasswordOptions = args.constraints[0] || {};
          const requirements: string[] = [];

          const minLength = opts.minLength ?? 8;
          requirements.push(`at least ${minLength} characters`);

          if (opts.requireUppercase ?? true) {
            requirements.push('one uppercase letter');
          }

          if (opts.requireLowercase ?? true) {
            requirements.push('one lowercase letter');
          }

          if (opts.requireNumbers ?? true) {
            requirements.push('one number');
          }

          if (opts.requireSpecialChars ?? true) {
            requirements.push(
              'one special character (@$!%*?&#^()_+-=[]{};\':"\\|,.<>/)',
            );
          }

          if (!(opts.allowWhitespace ?? false)) {
            requirements.push('no whitespace');
          }

          return `Password must contain ${requirements.join(', ')}`;
        },
      },
    });
  };
}
