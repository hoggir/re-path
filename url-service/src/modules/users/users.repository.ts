import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from './entities/user.entity';
import { TypeOrmBaseRepository } from '../../database/repositories/typeorm-base.repository';

@Injectable()
export class UsersRepository extends TypeOrmBaseRepository<User> {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
  ) {
    super(userRepository);
  }

  async findByEmail(email: string): Promise<User | null> {
    return this.userRepository.findOne({
      where: { email: email.toLowerCase(), isDeleted: false },
      select: {
        id: true,
        email: true,
        name: true,
        password: true,
        role: true,
        isActive: true,
        isEmailVerified: true,
        refreshToken: true,
        lastLoginAt: true,
        emailVerificationToken: true,
        emailVerificationExpires: true,
        passwordResetToken: true,
        passwordResetExpires: true,
        userCode: true,
        isDeleted: true,
        createdAt: true,
        updatedAt: true,
        deletedAt: true,
      },
    });
  }

  async findById(id: number): Promise<User | null> {
    return this.userRepository.findOne({
      where: { id, isDeleted: false },
      select: {
        id: true,
        email: true,
        name: true,
        role: true,
        isActive: true,
        isEmailVerified: true,
        lastLoginAt: true,
        userCode: true,
        isDeleted: true,
        createdAt: true,
        updatedAt: true,
      },
    });
  }

  async findByEmailWithPassword(email: string): Promise<User | null> {
    return this.userRepository.findOne({
      where: { email: email.toLowerCase(), isDeleted: false },
      select: {
        id: true,
        email: true,
        name: true,
        password: true,
        role: true,
        isActive: true,
        isEmailVerified: true,
        lastLoginAt: true,
        userCode: true,
        isDeleted: true,
        createdAt: true,
        updatedAt: true,
      },
    });
  }

  async updateRefreshToken(
    userId: number,
    refreshToken: string | null,
  ): Promise<void> {
    await this.userRepository.update(
      { id: userId },
      { refreshToken: refreshToken as any },
    );
  }

  async updateLastLogin(userId: number): Promise<void> {
    await this.userRepository.update(
      { id: userId },
      { lastLoginAt: new Date() },
    );
  }

  async findByRefreshToken(refreshToken: string): Promise<User | null> {
    return this.userRepository.findOne({
      where: { refreshToken, isDeleted: false },
      select: {
        id: true,
        email: true,
        name: true,
        role: true,
        isActive: true,
        isEmailVerified: true,
        refreshToken: true,
        lastLoginAt: true,
        userCode: true,
        isDeleted: true,
        createdAt: true,
        updatedAt: true,
      },
    });
  }
}
