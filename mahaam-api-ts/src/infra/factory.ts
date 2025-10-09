import { MiddlewareConsumer, Module, NestModule } from '@nestjs/common';
import { APP_FILTER, APP_INTERCEPTOR } from '@nestjs/core';
import { DefaultPlansService } from '../feat/plans/plans.service';
import { DefaultPlansRepo } from '../feat/plans/plans.repo';
import { DefaultPlanMembersRepo } from '../feat/plans/plan-members.repo';
import { DefaultPlansController } from '../feat/plans/plans.controller';
import { DefaultUsersRepo } from '../feat/users/users.repo';
import { DefaultSuggestedEmailRepo } from '../feat/users/suggested-email.repo';
import { DefaultTasksService } from 'src/feat/tasks/tasks.service';
import { DefaultTasksRepo } from 'src/feat/tasks/tasks.repo';
import { DefaultTasksController } from 'src/feat/tasks/tasks.controller';
import { DefaultUsersService } from 'src/feat/users/users.service';
import { DefaultDeviceRepo } from 'src/feat/users/device.repo';
import { DefaultUsersController } from 'src/feat/users/users.controller';
import { DefaultHealthService } from './monitor/health.service';
import { DefaultHealthRepo } from './monitor/health.repo';
import { DefaultTrafficRepo } from './monitor/traffic.repo';
import { DefaultLogRepo } from './monitor/log.repo';
import DefaultLog from './log';
import { Starter } from './starter';
import { DefaultAuditController } from './monitor/audit.controller';
import { DefaultHealthController } from './monitor/health.controller';
import { AppMiddleware, TrafficInterceptor } from './middleware';

@Module({
  imports: [],
  controllers: [
    DefaultUsersController,
    DefaultPlansController,
    DefaultTasksController,
    DefaultHealthController,
    DefaultAuditController,
  ],
  providers: [
    { provide: 'PlansRepo', useClass: DefaultPlansRepo },
    { provide: 'PlanMembersRepo', useClass: DefaultPlanMembersRepo },
    { provide: 'UsersRepo', useClass: DefaultUsersRepo },
    { provide: 'SuggestedEmailRepo', useClass: DefaultSuggestedEmailRepo },
    { provide: 'TasksRepo', useClass: DefaultTasksRepo },
    { provide: 'DeviceRepo', useClass: DefaultDeviceRepo },
    { provide: 'HealthRepo', useClass: DefaultHealthRepo },
    { provide: 'TrafficRepo', useClass: DefaultTrafficRepo },
    { provide: 'LogRepo', useClass: DefaultLogRepo },
    { provide: 'PlansService', useClass: DefaultPlansService },
    { provide: 'TasksService', useClass: DefaultTasksService },
    { provide: 'UsersService', useClass: DefaultUsersService },
    { provide: 'HealthService', useClass: DefaultHealthService },
    { provide: 'Starter', useClass: Starter },
    { provide: 'Log', useClass: DefaultLog },
    { provide: APP_INTERCEPTOR, useClass: TrafficInterceptor },
  ],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(AppMiddleware).exclude('/health', '/swagger', '/api-docs').forRoutes('*path');
  }
}
