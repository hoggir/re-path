import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { BaseSchema } from '../../../database/schemas/base.schema';

@Schema({ collection: 'urls', timestamps: true })
export class Url extends BaseSchema {
  @Prop({ required: true, unique: true })
  shortCode: string;

  @Prop({ required: true })
  longUrl: string;

  @Prop({ required: true })
  activeUntil: Date;

  @Prop({ default: true })
  isActive: boolean;

  @Prop({ default: 0 })
  accessCount: number;
}

export const UrlSchema = SchemaFactory.createForClass(Url);

// Indexes
UrlSchema.index({ shortCode: 1 }, { unique: true });
UrlSchema.index({ createdAt: -1 });
UrlSchema.index({ isDeleted: 1, isActive: 1 });

// Transform to JSON
UrlSchema.set('toJSON', {
  virtuals: true,
  transform: (doc, ret) => {
    const { __v, _id, ...remaining } = ret;
    return remaining;
  },
});
