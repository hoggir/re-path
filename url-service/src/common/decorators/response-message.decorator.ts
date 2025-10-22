import { SetMetadata } from '@nestjs/common';

export const RESPONSE_MESSAGE_KEY = 'response_message';

/**
 * Decorator to set custom response message for successful responses
 * @param message - The message to be included in the response
 * @example
 * @ResponseMessage('User created successfully')
 * async createUser() { ... }
 */
export const ResponseMessage = (message: string) =>
  SetMetadata(RESPONSE_MESSAGE_KEY, message);
