import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  UpdateDateColumn,
  DeleteDateColumn,
  Index,
  BeforeInsert,
  BeforeUpdate,
  AfterLoad,
} from 'typeorm';

@Entity('users')
@Index(['email'], { unique: true })
@Index(['createdAt'])
@Index(['isDeleted', 'isActive'])
@Index(['emailVerificationToken'])
@Index(['passwordResetToken'])
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  encryptedId?: string;

  @Column({ type: 'varchar', length: 255, unique: true })
  email: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'varchar', length: 255, select: false })
  password: string;

  @Column({
    type: 'enum',
    enum: ['user', 'admin', 'moderator'],
    default: 'user',
  })
  role: string;

  @Column({ type: 'boolean', default: true })
  isActive: boolean;

  @Column({ type: 'boolean', default: false })
  isEmailVerified: boolean;

  @Column({ type: 'text', nullable: true, select: false })
  refreshToken?: string;

  @Column({ type: 'timestamp', nullable: true })
  lastLoginAt?: Date;

  @Column({ type: 'varchar', length: 255, nullable: true })
  emailVerificationToken?: string;

  @Column({ type: 'timestamp', nullable: true })
  emailVerificationExpires?: Date;

  @Column({ type: 'varchar', length: 255, nullable: true })
  passwordResetToken?: string;

  @Column({ type: 'timestamp', nullable: true })
  passwordResetExpires?: Date;

  @Column({ type: 'int', unique: true, generated: 'increment' })
  userCode: number;

  @Column({ type: 'boolean', default: false })
  isDeleted: boolean;

  @CreateDateColumn({ type: 'timestamp' })
  createdAt: Date;

  @UpdateDateColumn({ type: 'timestamp' })
  updatedAt: Date;

  @DeleteDateColumn({ type: 'timestamp', nullable: true })
  deletedAt?: Date;

  @BeforeInsert()
  @BeforeUpdate()
  emailToLowerCase() {
    if (this.email) {
      this.email = this.email.toLowerCase().trim();
    }
  }
}
