import { Global, Module } from '@nestjs/common';
import { IdEncryptionService } from './utils/id-encryption.service';

@Global()
@Module({
  providers: [IdEncryptionService],
  exports: [IdEncryptionService],
})
export class CommonModule {}
