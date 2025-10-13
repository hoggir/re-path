import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import Hashids from 'hashids';

@Injectable()
export class IdEncryptionService {
  private hashids: Hashids;

  constructor(private readonly configService: ConfigService) {
    const salt = this.configService.get<string>('ID_ENCRYPTION_SALT') || 'default-salt-change-in-production';
    const minLength = 8;
    this.hashids = new Hashids(salt, minLength);
  }

  /**
   * Encrypt numeric ID to string
   */
  encryptId(id: number): string {
    return this.hashids.encode(id);
  }

  /**
   * Decrypt string ID to number
   */
  decryptId(encryptedId: string): number | null {
    const decoded = this.hashids.decode(encryptedId);
    return decoded.length > 0 ? Number(decoded[0]) : null;
  }

  /**
   * Encrypt multiple IDs
   */
  encryptIds(ids: number[]): string[] {
    return ids.map(id => this.encryptId(id));
  }

  /**
   * Decrypt multiple IDs
   */
  decryptIds(encryptedIds: string[]): number[] {
    return encryptedIds
      .map(id => this.decryptId(id))
      .filter((id): id is number => id !== null);
  }
}
