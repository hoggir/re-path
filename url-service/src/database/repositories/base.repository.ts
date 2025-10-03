import {
  Document,
  FilterQuery,
  Model,
  UpdateQuery,
  QueryOptions,
} from 'mongoose';
import { NotFoundException } from '../../common/exceptions/custom-exceptions';

export abstract class BaseRepository<T extends Document> {
  constructor(protected readonly model: Model<T>) {}

  async create(createDto: Partial<T>): Promise<T> {
    const created = new this.model(createDto);
    return created.save();
  }

  async findAll(
    filter: FilterQuery<T> = {},
    options: QueryOptions = {},
  ): Promise<T[]> {
    return this.model
      .find({ ...filter, isDeleted: false })
      .setOptions(options)
      .exec();
  }

  async findOne(
    filter: FilterQuery<T>,
    options: QueryOptions = {},
  ): Promise<T | null> {
    return this.model
      .findOne({ ...filter, isDeleted: false })
      .setOptions(options)
      .exec();
  }

  async findById(id: string, options: QueryOptions = {}): Promise<T | null> {
    return this.model
      .findOne({ _id: id, isDeleted: false } as FilterQuery<T>)
      .setOptions(options)
      .exec();
  }

  async findByIdOrFail(id: string, options: QueryOptions = {}): Promise<T> {
    const document = await this.findById(id, options);
    if (!document) {
      throw new NotFoundException(
        `${this.model.modelName} with ID ${id} not found`,
      );
    }
    return document;
  }

  async update(id: string, updateDto: UpdateQuery<T>): Promise<T | null> {
    return this.model
      .findOneAndUpdate(
        { _id: id, isDeleted: false } as FilterQuery<T>,
        { $set: updateDto },
        { new: true },
      )
      .exec();
  }

  async delete(id: string): Promise<T | null> {
    return this.model.findByIdAndDelete(id).exec();
  }

  // Soft delete
  async softDelete(id: string): Promise<T | null> {
    return this.model
      .findByIdAndUpdate(
        id,
        {
          $set: {
            isDeleted: true,
            deletedAt: new Date(),
          },
        },
        { new: true },
      )
      .exec();
  }

  async restore(id: string): Promise<T | null> {
    return this.model
      .findByIdAndUpdate(
        id,
        {
          $set: {
            isDeleted: false,
            deletedAt: null,
          },
        },
        { new: true },
      )
      .exec();
  }

  async count(filter: FilterQuery<T> = {}): Promise<number> {
    return this.model.countDocuments({ ...filter, isDeleted: false }).exec();
  }

  async exists(filter: FilterQuery<T>): Promise<boolean> {
    const count = await this.model
      .countDocuments({ ...filter, isDeleted: false })
      .limit(1)
      .exec();
    return count > 0;
  }

  async paginate(
    filter: FilterQuery<T> = {},
    page: number = 1,
    limit: number = 10,
    sort: any = { createdAt: -1 },
  ) {
    const skip = (page - 1) * limit;

    const [data, total] = await Promise.all([
      this.model
        .find({ ...filter, isDeleted: false })
        .sort(sort)
        .skip(skip)
        .limit(limit)
        .exec(),
      this.count(filter),
    ]);

    return {
      data,
      pagination: {
        page,
        limit,
        total,
        totalPages: Math.ceil(total / limit),
        hasNext: page * limit < total,
        hasPrev: page > 1,
      },
    };
  }
}
