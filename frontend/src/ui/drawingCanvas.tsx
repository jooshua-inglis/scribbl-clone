import { MouseEvent, useCallback, useEffect, useRef, useState } from 'react';
import { produce } from 'immer'

export type Point = {
    x: number,
    y: number,
}
export type Line = {
    points: Point[]
    size: number
    rgb: [number, number, number]
}

function getDistance(p1: { x: number; y: number }, p2: { x: number; y: number }) {
    return Math.sqrt((p1.x - p2.x) ** 2 + (p1.y - p2.y) ** 2);
}

function getMousePos(canvas: HTMLCanvasElement, evt: MouseEvent) {
    const rect = canvas.getBoundingClientRect();
    return {
        x: evt.clientX - rect.left,
        y: evt.clientY - rect.top,
    };
}


export function DrawingCanvas({
    lines = [],
    width,
    height,
    onDraw,
    drawable = true,
}: {
    lines: Line[];
    width: number;
    height: number;
    onDraw?: (lines: Line[]) => unknown;
    drawable?: boolean;
}) {
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const [lastMarkPos, setLastMarkPos] = useState<{ x: number; y: number } | undefined>(undefined);
    const [drawing, setDrawing] = useState(false)
    const [drawingLineIndex, setDrawingLineIndex] = useState<number | undefined>(undefined)

    useEffect(() => {
        const canvas = canvasRef.current;
        if (!canvas) {
            return;
        }
        const ctx = canvas.getContext('2d');
        if (!ctx) {
            return;
        }
        const interval = setInterval(() => {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            ctx.beginPath();
            for (const line of lines) {
                ctx.lineWidth = line.size
                ctx.strokeStyle = `rgb(${line.rgb[0]} ${line.rgb[1]} ${line.rgb[2]})`
                if(line.points.length === 0) {
                    continue
                }
                ctx.moveTo(line.points[0].x, line.points[0].y);
                for (const point of line.points) {
                    ctx.lineTo(point.x, point.y);
                }
            }
            ctx.stroke();
        }, 1)
        return () => clearInterval(interval)
    }, [lines])



    const onDrag = useCallback(
        (e: MouseEvent<HTMLCanvasElement>) => {
            const mousePos = getMousePos(e.currentTarget, e)
            // Only draw new point if it's a set distance away from the last mouse pos. This will
            // help with debouncing
            if (lastMarkPos && getDistance(lastMarkPos, mousePos) < 0.2) {
                return
            }
            if (!drawing) {
                return
            }
            if(drawingLineIndex == null) {
                return
            }

            if (onDraw) {
                const nextLines = produce(lines, draftLines => {
                    draftLines[drawingLineIndex].points.push(mousePos)
                })
                onDraw(nextLines)
            }
            setLastMarkPos(mousePos)
            setDrawing(true)
        },
        [lastMarkPos, drawing, drawingLineIndex, onDraw, lines]
    );

    const startDrawing = useCallback(() => {
        if (!drawable) {
            return
        }
        setDrawing(true)

        if (onDraw) {
            const newLine: Line = {
                points: [],
                rgb: [20, 200, 20],
                size: 5
            }
            onDraw([...lines, newLine])
        }
        setDrawingLineIndex(lines.length)
    }, [lines, onDraw, drawable])


    const stopDrawing = useCallback(() => {
        setLastMarkPos(undefined)
        setDrawingLineIndex(undefined)
        setDrawing(false)
    }, [])

    return (
        <canvas
            onMouseDown={startDrawing}
            onMouseUp={stopDrawing}
            onMouseLeave={stopDrawing}
            onMouseMove={onDrag}
            style={{ backgroundColor: 'white' }}
            height={height}
            width={width}
            ref={canvasRef}
        />
    );
}
