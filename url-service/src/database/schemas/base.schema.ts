import { Prop, Schema } from '@nestjs/mongoose';
import { Document, Types } from 'mongoose';

@Schema({ timestamps: true })
export class BaseSchema extends Document {
   @Prop({ type: Types.ObjectId, auto: true })
   declare _id: Types.ObjectId;

  @Prop({ default: Date.now })
  createdAt: Date;

  @Prop({ default: Date.now })
  updatedAt: Date;

  @Prop({ default: null })
  deletedAt?: Date;

  @Prop({ default: false })
  isDeleted: boolean;
}