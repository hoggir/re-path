import { Injectable, UnauthorizedException } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';
import { ConfigService } from '@nestjs/config';
import { RedisService } from '../../../cache/redis.service';
import { AuthService } from '../auth.service';
import { IdEncryptionService } from '../../../common/utils/id-encryption.service';

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(
    private readonly configService: ConfigService,
    private readonly authService: AuthService,
    private readonly redisService: RedisService,
    private readonly idEncryptionService: IdEncryptionService,
  ) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: configService.get<string>('auth.jwt.secret') as string,
    });
  }

  async validate(payload: any) {
    const cacheKey = `user:${payload.sub}`;
    let user = await this.redisService.get(cacheKey);

    if (!user) {
      console.log('💾 Cache miss - querying database...');
      user = await this.authService.validateUser(payload.sub);

      if (!user) {
        throw new UnauthorizedException('User not found or inactive');
      }

      await this.redisService.set(cacheKey, user, 900);
    } else {
      console.log('⚡ Cache hit - using cached data');
    }

    return {
      userId: payload.sub,
      encryptedUserId: this.idEncryptionService.encryptId(payload.sub),
      email: payload.email,
      role: payload.role,
    };
  }
}
