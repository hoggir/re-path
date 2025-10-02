import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [ConfigModule.forRoot({
          // Optional: configure path to .env, make global, etc.
        }),],
  controllers: [],
  providers: [],
})
export class AppModule {}
