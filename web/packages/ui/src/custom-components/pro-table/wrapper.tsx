import {
  DndContext,
  DragEndEvent,
  KeyboardSensor,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';

export function ProTableWrapper<TData extends { id?: string }>({
  children,
  onSort,
  data,
  setData,
}: {
  children: React.ReactNode;
  onSort?: (
    sourceId: string | number,
    targetId: string | number | null,
    items: TData[],
  ) => Promise<TData[]>;
  data: TData[];
  setData: React.Dispatch<React.SetStateAction<TData[]>>;
}) {
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    if (onSort) {
      const updatedData = await onSort(active.id, over?.id || null, data);
      setData(updatedData);
    }
  };
  if (!onSort) return children;
  return (
    <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
      <SortableContext
        items={data.map((item) => String(item.id))}
        strategy={verticalListSortingStrategy}
      >
        {children}
      </SortableContext>
    </DndContext>
  );
}
