import {
  Repository,
  FindOptionsWhere,
  FindManyOptions,
  DeepPartial,
  FindOptionsOrder,
} from 'typeorm';
import { NotFoundException } from '@nestjs/common';

export abstract class TypeOrmBaseRepository<
  T extends { id: number; isDeleted: boolean },
> {
  constructor(protected readonly repository: Repository<T>) {}

  async create(createDto: DeepPartial<T>): Promise<T> {
    const entity = this.repository.create(createDto);
    return this.repository.save(entity);
  }

  async findAll(
    filter: FindOptionsWhere<T> = {} as FindOptionsWhere<T>,
    options: FindManyOptions<T> = {},
  ): Promise<T[]> {
    return this.repository.find({
      where: { ...filter, isDeleted: false } as FindOptionsWhere<T>,
      ...options,
    });
  }

  async findOne(
    filter: FindOptionsWhere<T>,
    options: FindManyOptions<T> = {},
  ): Promise<T | null> {
    return this.repository.findOne({
      where: { ...filter, isDeleted: false } as FindOptionsWhere<T>,
      ...options,
    });
  }

  async findById(
    id: number,
    options: FindManyOptions<T> = {},
  ): Promise<T | null> {
    return this.repository.findOne({
      where: { id, isDeleted: false } as FindOptionsWhere<T>,
      ...options,
    });
  }

  async findByIdOrFail(
    id: number,
    options: FindManyOptions<T> = {},
  ): Promise<T> {
    const entity = await this.findById(id, options);
    if (!entity) {
      throw new NotFoundException(
        `${this.repository.metadata.name} with ID ${id} not found`,
      );
    }
    return entity;
  }

  async update(id: number, updateDto: DeepPartial<T>): Promise<T | null> {
    await this.repository.update(
      { id, isDeleted: false } as FindOptionsWhere<T>,
      updateDto as any,
    );
    return this.findById(id);
  }

  async delete(id: number): Promise<void> {
    await this.repository.delete(id);
  }

  async softDelete(id: number): Promise<T | null> {
    await this.repository.update(id, {
      isDeleted: true,
      deletedAt: new Date(),
    } as any);
    return this.repository.findOne({ where: { id } as FindOptionsWhere<T> });
  }

  async restore(id: number): Promise<T | null> {
    await this.repository.update(id, {
      isDeleted: false,
      deletedAt: null,
    } as any);
    return this.findById(id);
  }

  async count(
    filter: FindOptionsWhere<T> = {} as FindOptionsWhere<T>,
  ): Promise<number> {
    return this.repository.count({
      where: { ...filter, isDeleted: false } as FindOptionsWhere<T>,
    });
  }

  async exists(filter: FindOptionsWhere<T>): Promise<boolean> {
    const count = await this.repository.count({
      where: { ...filter, isDeleted: false } as FindOptionsWhere<T>,
      take: 1,
    });
    return count > 0;
  }

  async paginate(
    filter: FindOptionsWhere<T> = {} as FindOptionsWhere<T>,
    page: number = 1,
    limit: number = 10,
    sort?: FindOptionsOrder<T>,
  ) {
    const skip = (page - 1) * limit;

    const [data, total] = await Promise.all([
      this.repository.find({
        where: { ...filter, isDeleted: false } as FindOptionsWhere<T>,
        order: sort || ({ createdAt: 'DESC' } as any),
        skip,
        take: limit,
      }),
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
