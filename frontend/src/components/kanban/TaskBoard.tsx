import { useEffect, useMemo, useState } from 'react'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import {
  Check,
  Circle,
  Flame,
  GripVertical,
  LogOut,
  Plus,
  RefreshCcw,
  Timer,
  UserCircle,
  ChevronDown,
} from 'lucide-react'

import { useAuth } from '../../integrations/auth/auth-context'

import { Button } from '../ui/button'
import { Checkbox } from '../ui/checkbox'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu'
import { Input } from '../ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../ui/select'
import { ScrollArea } from '../ui/scroll-area'
import { Tooltip, TooltipContent, TooltipTrigger } from '../ui/tooltip'
import { cn } from '../../lib/utils'
import { createTask, listTasks, updateTask } from '../../lib/tasks'
import type { SubmitEvent } from 'react'
import type { DragStart, DropResult } from '@hello-pangea/dnd'
import type {
  CreateTaskInput,
  Task,
  TaskPriority,
  TaskStatus,
  UpdateTaskInput,
} from '../../lib/types'

const columnMeta: Array<{
  id: TaskStatus
  title: string
  tone: string
  icon: typeof Circle
  hint: string
}> = [
  {
    id: 'OPEN',
    title: 'Backlog',
    tone: 'text-amber-300',
    icon: Circle,
    hint: 'Queued for attention',
  },
  {
    id: 'IN_PROGRESS',
    title: 'In Progress',
    tone: 'text-amber-200',
    icon: Timer,
    hint: 'Actively being built',
  },
  {
    id: 'DONE',
    title: 'Done',
    tone: 'text-amber-100',
    icon: Check,
    hint: 'Shipped and verified',
  },
]

const priorityTone: Record<TaskPriority, string> = {
  HIGH: 'bg-amber-400/15 text-amber-200 border-amber-400/40',
  MEDIUM: 'bg-amber-200/10 text-amber-100 border-amber-200/30',
  LOW: 'bg-slate-200/10 text-slate-200/70 border-slate-200/20',
}

const emptyOrder: Record<TaskStatus, Array<string>> = {
  OPEN: [],
  IN_PROGRESS: [],
  DONE: [],
}

export function TaskBoard() {
  const navigate = useNavigate()
  const { user, signOut } = useAuth()
  const queryClient = useQueryClient()
  const [title, setTitle] = useState('')
  const [priority, setPriority] = useState<TaskPriority>('MEDIUM')
  const [orderByStatus, setOrderByStatus] = useState(emptyOrder)
  const [activeDragId, setActiveDragId] = useState<string | null>(null)

  const tasksQuery = useQuery({
    queryKey: ['tasks'],
    queryFn: listTasks,
  })

  const tasks = tasksQuery.data ?? []
  const tasksById = useMemo(
    () => new Map(tasks.map((task) => [task.id, task])),
    [tasks],
  )

  useEffect(() => {
    if (!tasks.length) {
      setOrderByStatus(emptyOrder)
      return
    }

    setOrderByStatus((prev) => {
      const next: Record<TaskStatus, Array<string>> = { ...prev }

      for (const column of columnMeta) {
        const incoming = tasks
          .filter((task) => task.status === column.id)
          .map((task) => task.id)
        const preserved = prev[column.id].filter((id) => incoming.includes(id))
        const additions = incoming.filter((id) => !preserved.includes(id))
        next[column.id] = [...preserved, ...additions]
      }

      return next
    })
  }, [tasks])

  const createMutation = useMutation({
    mutationFn: (input: CreateTaskInput) => createTask(input),
    onMutate: async (input) => {
      await queryClient.cancelQueries({ queryKey: ['tasks'] })
      const previous = queryClient.getQueryData<Array<Task>>(['tasks']) ?? []
      const tempId = `temp-${Date.now()}`

      const optimisticTask: Task = {
        id: tempId,
        title: input.title,
        status: 'OPEN',
        priority: input.priority,
        assignee_id: input.assignee_id ?? null,
        created_by: 'you',
        created_at: new Date().toISOString(),
      }

      queryClient.setQueryData<Array<Task>>(
        ['tasks'],
        [optimisticTask, ...previous],
      )
      setOrderByStatus((prev) => ({
        ...prev,
        OPEN: [tempId, ...prev.OPEN],
      }))

      return { previous, tempId }
    },
    onError: (_error, _input, context) => {
      if (context?.previous) {
        queryClient.setQueryData(['tasks'], context.previous)
      }
      toast.error('Task create failed. Please try again.')
    },
    onSuccess: (created, _input, context) => {
      queryClient.setQueryData<Array<Task>>(['tasks'], (current = []) =>
        current.map((task) => (task.id === context.tempId ? created : task)),
      )
      setOrderByStatus((prev) => ({
        ...prev,
        OPEN: prev.OPEN.map((id) => (id === context.tempId ? created.id : id)),
      }))
      toast.success('Task created')
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateTaskInput }) =>
      updateTask(id, input),
    onMutate: async ({ id, input }) => {
      await queryClient.cancelQueries({ queryKey: ['tasks'] })
      const previous = queryClient.getQueryData<Array<Task>>(['tasks']) ?? []

      queryClient.setQueryData<Array<Task>>(['tasks'], (current = []) =>
        current.map((task) =>
          task.id === id
            ? {
                ...task,
                ...input,
              }
            : task,
        ),
      )

      if (input.status) {
        setOrderByStatus((prev) => {
          const next: Record<TaskStatus, Array<string>> = {
            OPEN: prev.OPEN.filter((taskId) => taskId !== id),
            IN_PROGRESS: prev.IN_PROGRESS.filter((taskId) => taskId !== id),
            DONE: prev.DONE.filter((taskId) => taskId !== id),
          }
          next[input.status as TaskStatus] = [
            id,
            ...next[input.status as TaskStatus],
          ]
          return next
        })
      }

      return { previous }
    },
    onError: (_error, _input, context) => {
      if (context?.previous) {
        queryClient.setQueryData(['tasks'], context.previous)
      }
      toast.error('Task update failed. Please try again.')
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })

  const orderedTasks = useMemo(() => {
    const orderFor = (status: TaskStatus) =>
      orderByStatus[status].length
        ? orderByStatus[status]
            .map((id) => tasksById.get(id))
            .filter((task): task is Task => Boolean(task))
        : tasks.filter((task) => task.status === status)

    return {
      OPEN: orderFor('OPEN'),
      IN_PROGRESS: orderFor('IN_PROGRESS'),
      DONE: orderFor('DONE'),
    }
  }, [orderByStatus, tasks, tasksById])

  const handleCreate = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault()
    const trimmed = title.trim()
    if (!trimmed || createMutation.isPending) return

    createMutation.mutate({
      title: trimmed,
      priority,
      description: '',
    })

    setTitle('')
  }

  const handleDragStart = (start: DragStart) => {
    setActiveDragId(start.draggableId)
  }

  const handleDragEnd = (result: DropResult) => {
    const { destination, source, draggableId } = result
    setActiveDragId(null)
    if (!destination) return

    const sourceStatus = source.droppableId as TaskStatus
    const destinationStatus = destination.droppableId as TaskStatus

    setOrderByStatus((prev) => {
      const next = { ...prev }
      const sourceOrder = Array.from(next[sourceStatus])
      const [removed] = sourceOrder.splice(source.index, 1)

      if (sourceStatus === destinationStatus) {
        sourceOrder.splice(destination.index, 0, removed)
        next[sourceStatus] = sourceOrder
        return next
      }

      const destinationOrder = Array.from(next[destinationStatus])
      destinationOrder.splice(destination.index, 0, removed)
      next[sourceStatus] = sourceOrder
      next[destinationStatus] = destinationOrder
      return next
    })

    if (sourceStatus !== destinationStatus) {
      updateMutation.mutate({
        id: draggableId,
        input: { status: destinationStatus },
      })
    }
  }

  const resetOrder = () => {
    setOrderByStatus(() => {
      const next: Record<TaskStatus, Array<string>> = { ...emptyOrder }
      for (const column of columnMeta) {
        next[column.id] = tasks
          .filter((task) => task.status === column.id)
          .map((task) => task.id)
      }
      return next
    })
    toast.success('Order reset')
  }

  const handleSignOut = () => {
    signOut()
    toast.success('Signed out')
    navigate({ to: '/login', replace: true })
  }

  const isInitialLoading = tasksQuery.isLoading
  const isRefreshing = tasksQuery.isFetching && !tasksQuery.isLoading

  return (
    <div className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-8 px-6 pb-16 pt-12">
      <header className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-3 text-sm uppercase tracking-[0.32em] text-amber-200/70">
            <Flame className="h-4 w-4" />
            Dispatch Board
          </div>
          <h1 className="mt-2 text-3xl font-semibold text-amber-100">
            Serverless Kanban
          </h1>
          <p className="mt-1 text-sm text-slate-200/60">
            Optimistic task flow with cold-start aware loading.
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-3">
          <Button
            asChild
            variant="outline"
            className={cn(
              'rounded-full border-white/10 bg-white/5 px-4 text-xs font-medium uppercase tracking-[0.2em] text-amber-100/80 hover:bg-white/10',
              isRefreshing && 'opacity-60',
            )}
          >
            <motion.button
              type="button"
              onClick={() => tasksQuery.refetch()}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              className="flex items-center gap-2"
            >
              <RefreshCcw
                className={cn('h-4 w-4', isRefreshing && 'animate-spin')}
              />
              Sync
            </motion.button>
          </Button>
          <div className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-xs uppercase tracking-[0.2em] text-amber-100/70">
            {tasks.length} tasks
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="outline"
                className="flex items-center gap-2 rounded-full border-white/10 bg-white/5 px-4 text-xs font-medium uppercase tracking-[0.18em] text-amber-100/80 hover:bg-white/10"
              >
                <UserCircle className="h-4 w-4" />
                Account
                <ChevronDown className="h-3.5 w-3.5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-52">
              <DropdownMenuLabel className="text-xs uppercase tracking-[0.2em] text-amber-100/60">
                {user?.email ?? 'Signed in'}
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => tasksQuery.refetch()}>
                <RefreshCcw className="h-4 w-4" />
                Refresh tasks
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={resetOrder}>
                Reset column order
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={handleSignOut}
                className="text-red-400 focus:bg-red-500/10"
              >
                <LogOut className="h-4 w-4" />
                Sign out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </header>

      <form
        onSubmit={handleCreate}
        className="flex flex-wrap items-center gap-3 rounded-2xl border border-white/10 bg-white/5 p-4"
      >
        <div className="flex min-w-55 flex-1 items-center gap-3">
          <span className="text-xs uppercase tracking-[0.2em] text-amber-100/60">
            New task
          </span>
          <Input
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            placeholder="Define the next move"
            className="h-10 border-white/10 bg-black/30 text-sm text-amber-100 placeholder:text-amber-100/30 focus-visible:ring-amber-200/20"
          />
        </div>
        <Select
          value={priority}
          onValueChange={(value) => setPriority(value as TaskPriority)}
        >
          <SelectTrigger className="rounded-full border-white/10 bg-black/30 text-xs uppercase tracking-[0.2em] text-amber-100/80">
            <SelectValue placeholder="Priority" />
          </SelectTrigger>
          <SelectContent className="border-white/10 bg-black/90 text-amber-50">
            <SelectItem value="LOW">Low</SelectItem>
            <SelectItem value="MEDIUM">Medium</SelectItem>
            <SelectItem value="HIGH">High</SelectItem>
          </SelectContent>
        </Select>
        <Button
          asChild
          className="rounded-full bg-amber-200/90 px-4 text-xs font-semibold uppercase tracking-[0.2em] text-black hover:bg-amber-200"
        >
          <motion.button
            type="submit"
            whileHover={{ scale: 1.03 }}
            whileTap={{ scale: 0.97 }}
            className="flex items-center gap-2"
          >
            <Plus className="h-4 w-4" />
            Create
          </motion.button>
        </Button>
      </form>

      {tasksQuery.isError && (
        <div className="rounded-2xl border border-amber-400/30 bg-amber-400/10 px-4 py-3 text-sm text-amber-100">
          Task sync failed. Retrying is enabled for transient gateway errors.
        </div>
      )}

      <DragDropContext onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
        <section className="grid gap-6 md:grid-cols-3">
          {columnMeta.map((column) => {
            const ColumnIcon = column.icon
            const columnTasks = orderedTasks[column.id]

            return (
              <Droppable key={column.id} droppableId={column.id}>
                {(provided, snapshot) => (
                  <div
                    className={cn(
                      'flex h-[70vh] min-h-105 flex-col gap-4 overflow-hidden rounded-3xl border border-white/10 bg-white/5 p-4 transition duration-100',
                      snapshot.isDraggingOver &&
                        'border-amber-200/40 bg-amber-200/10',
                    )}
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <div className="flex items-center gap-2 text-sm text-amber-100">
                          <ColumnIcon className={cn('h-4 w-4', column.tone)} />
                          <span className="font-medium tracking-widest">
                            {column.title}
                          </span>
                        </div>
                        <p className="mt-1 text-xs text-slate-200/50">
                          {column.hint}
                        </p>
                      </div>
                      <div className="rounded-full border border-white/10 bg-black/30 px-3 py-1 text-xs text-amber-100/70">
                        {columnTasks.length}
                      </div>
                    </div>

                    <ScrollArea
                      className="scroll-area-minimal min-h-0 flex-1"
                      viewportClassName="flex min-h-0 flex-col gap-3 pb-6 pr-2"
                      viewportProps={{
                        ref: provided.innerRef,
                        ...provided.droppableProps,
                      }}
                    >
                      {isInitialLoading
                        ? Array.from({ length: 3 }).map((_, index) => (
                            <div
                              key={`${column.id}-skeleton-${index}`}
                              className="h-24 animate-pulse rounded-2xl border border-white/10 bg-white/5"
                            />
                          ))
                        : columnTasks.map((task, index) => {
                            const isActive = activeDragId === task.id
                            const shouldBlur =
                              Boolean(activeDragId) && !isActive

                            return (
                              <Draggable
                                key={task.id}
                                draggableId={task.id}
                                index={index}
                              >
                                {(dragProvided, dragSnapshot) => (
                                  <motion.article
                                    ref={dragProvided.innerRef}
                                    {...dragProvided.draggableProps}
                                    {...dragProvided.dragHandleProps}
                                    layout
                                    whileHover={{ y: -2 }}
                                    animate={{
                                      filter: shouldBlur
                                        ? 'blur(2px) saturate(0.85)'
                                        : 'blur(0px) saturate(1)',
                                      opacity: shouldBlur ? 0.8 : 1,
                                      scale: isActive ? 1.02 : 1,
                                    }}
                                    transition={{
                                      duration: 0.12,
                                      ease: 'easeOut',
                                    }}
                                    style={{ willChange: 'transform, filter' }}
                                    className={cn(
                                      'group transform-gpu rounded-2xl border border-white/10 bg-black/40 p-4 shadow-[0_10px_40px_rgba(10,10,14,0.35)] transition duration-100 hover:border-amber-200/40',
                                      dragSnapshot.isDragging &&
                                        'border-amber-200/60 bg-black/70',
                                    )}
                                  >
                                    <div className="flex items-start gap-3">
                                      <Tooltip>
                                        <TooltipTrigger asChild>
                                          <span
                                            aria-label="Drag task"
                                            className="mt-1 inline-flex h-8 w-8 items-center justify-center rounded-full border border-white/10 bg-white/5 text-amber-100/70 transition duration-100 hover:text-amber-100"
                                          >
                                            <GripVertical className="h-4 w-4" />
                                          </span>
                                        </TooltipTrigger>
                                        <TooltipContent
                                          side="top"
                                          className="border-white/10 bg-black/90 text-amber-100/80"
                                        >
                                          Drag to reorder
                                        </TooltipContent>
                                      </Tooltip>
                                      <div className="flex-1">
                                        <div className="flex items-center justify-between gap-2">
                                          <h3 className="text-sm font-medium text-amber-100">
                                            {task.title}
                                          </h3>
                                          <span
                                            className={cn(
                                              'rounded-full border px-2 py-0.5 text-[10px] uppercase tracking-[0.2em]',
                                              priorityTone[task.priority],
                                            )}
                                          >
                                            {task.priority}
                                          </span>
                                        </div>
                                        <div className="mt-2 flex items-center gap-3 text-xs text-slate-200/50">
                                          <span>Owner {task.created_by}</span>
                                          <span className="h-1 w-1 rounded-full bg-white/20" />
                                          <span>
                                            {new Date(
                                              task.created_at,
                                            ).toLocaleDateString()}
                                          </span>
                                        </div>
                                      </div>
                                      <Checkbox
                                        checked={task.status === 'DONE'}
                                        onCheckedChange={(checked) => {
                                          const nextStatus =
                                            checked === true ? 'DONE' : 'OPEN'
                                          if (nextStatus !== task.status) {
                                            updateMutation.mutate({
                                              id: task.id,
                                              input: { status: nextStatus },
                                            })
                                          }
                                        }}
                                        aria-label="Toggle done"
                                        className="size-7 rounded-full border-white/10 bg-white/5 text-amber-100 shadow-none transition duration-100 hover:border-amber-200/60"
                                      />
                                    </div>
                                  </motion.article>
                                )}
                              </Draggable>
                            )
                          })}
                      {provided.placeholder}
                    </ScrollArea>
                  </div>
                )}
              </Droppable>
            )
          })}
        </section>
      </DragDropContext>
    </div>
  )
}
