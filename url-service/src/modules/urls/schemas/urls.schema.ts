import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { BaseSchema } from '../../../database/schemas/base.schema';

@Schema({ collection: 'urls', timestamps: true })
export class Url extends BaseSchema {
  @Prop({ required: true, type: String })
  originalUrl: string;

  @Prop({ required: true, unique: true, sparse: true, type: String })
  shortCode: string;

  @Prop({ type: String })
  customAlias?: string;

  @Prop({ type: Number, required: true, index: true })
  userId: number;

  @Prop({ default: 0, type: Number })
  clickCount: number;

  @Prop({ type: Date })
  expiresAt?: Date;

  @Prop({ default: true, type: Boolean })
  isActive: boolean;

  @Prop({ type: String })
  title?: string;

  @Prop({ type: String })
  description?: string;

  // @Prop({
  //   type: [
  //     {
  //       clickedAt: { type: Date, default: Date.now },
  //       ipAddress: String,
  //       userAgent: String,
  //       referer: String,
  //       country: String,
  //       city: String,
  //     },
  //   ],
  //   default: [],
  // })
  // analytics: {
  //   clickedAt: Date;
  //   ipAddress: string;
  //   userAgent: string;
  //   referer: string;
  //   country?: string;
  //   city?: string;
  // }[];

  @Prop({ type: Object })
  metadata?: {
    domain?: string;
    protocol?: string;
    [key: string]: any;
  };
}

export const UrlSchema = SchemaFactory.createForClass(Url);

UrlSchema.index({ shortCode: 1, isDeleted: 1 }, { unique: true, sparse: true });
UrlSchema.index({ userId: 1, createdAt: -1 });
UrlSchema.index({ expiresAt: 1 }, { expireAfterSeconds: 0 });
UrlSchema.index({ isActive: 1 });

// Virtual for short URL
UrlSchema.virtual('shortUrl').get(function () {
  return `${process.env.BASE_URL || 'http://localhost:3010'}/${this.shortCode}`;
});

// Transform to JSON
UrlSchema.set('toJSON', {
  virtuals: true,
  transform: (doc, ret) => {
    const { __v, _id, ...remaining } = ret;
    return remaining;
  },
});
