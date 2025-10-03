import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { BaseSchema } from '../../../database/schemas/base.schema';
import { CounterDocument } from 'src/database/schemas/counter.schema';

@Schema({ collection: 'users', timestamps: true })
export class User extends BaseSchema {
  @Prop({ required: true, unique: true, trim: true, lowercase: true })
  email: string;

  @Prop({ required: true, trim: true })
  name: string;

  @Prop({ required: true, select: false })
  password: string;

  @Prop({ default: 'user', enum: ['user', 'admin', 'moderator'] })
  role: string;

  @Prop({ default: true })
  isActive: boolean;

  @Prop({ default: false })
  isEmailVerified: boolean;

  @Prop({ type: String, select: false })
  refreshToken?: string;

  @Prop({ type: Date })
  lastLoginAt?: Date;

  @Prop({ type: String })
  emailVerificationToken?: string;

  @Prop({ type: Date })
  emailVerificationExpires?: Date;

  @Prop({ type: String })
  passwordResetToken?: string;

  @Prop({ type: Date })
  passwordResetExpires?: Date;

  @Prop({ type: Number, unique: true })
  userCode: number;
}

export const UserSchema = SchemaFactory.createForClass(User);

// ...
UserSchema.pre<User>('save', async function (next) {
  if (this.isNew) {
    const counterModel = this.db.model<CounterDocument>('Counter');
    const counter = await counterModel.findOneAndUpdate(
      { id: 'userCode' }, // ⚠️ harus pakai "id", bukan "name"
      { $inc: { seq: 1 } },
      { new: true, upsert: true },
    );

    this.userCode = counter.seq;
  }
  next();
});

// Indexes
UserSchema.index({ email: 1 }, { unique: true });
UserSchema.index({ createdAt: -1 });
UserSchema.index({ isDeleted: 1, isActive: 1 });
UserSchema.index({ emailVerificationToken: 1 });
UserSchema.index({ passwordResetToken: 1 });

// Transform to JSON
UserSchema.set('toJSON', {
  virtuals: true,
  transform: (doc, ret) => {
    const {
      __v,
      _id,
      refreshToken,
      emailVerificationToken,
      passwordResetToken,
      ...remaining
    } = ret;
    return remaining;
  },
});
