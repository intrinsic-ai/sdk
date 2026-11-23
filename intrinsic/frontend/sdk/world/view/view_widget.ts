// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {Observable, Subject, Subscription} from 'rxjs';
import * as three from 'three';

import {IntrinsicCamera} from './camera';
import * as colors from './colors';
import {IntrOrbitControls} from './intr_orbit_controls';
import {createCircleSpriteMaterial} from './materials';
import {RendererClient, WorldRendererService} from './world_renderer_service';

const VIEW_WIDGET_BASE_PIXEL_SIZE = 65;
// Speed of animation. Will be scaled by delta time in update.
const ANIMATION_INTERPOLATION_RATE = 0.001;
// Controls accel/decel in the animation. Must be <= -1 where -1 is linear
// speed.
const ANIMATION_SMOOTHING_VALUE = -3.0;
const VIEW_WIDGET_POS_AXIS_SCALE = 0.6;
const VIEW_WIDGET_NEG_AXIS_SCALE = 0.45;

/**
 * The ViewWidget uses threeJS to render a widget in the scene that allows the
 * user to quickly determine current view direction, as well as align view
 * direction with any world axis. The widget's axis are always aligned with
 * the view direction of the main frontend camera. When the user clicks on
 * one of the widget's control points, quaternion interpolation is used to
 * blend the frontend camera to the new orientation.
 *
 * The widget renders into the given canvas until `dispose` is called. Its
 * render loop runs in whatever zone the widget is constructed in, so Angular
 * callers should construct it outside of the Angular zone.
 */
export class ViewWidget {
  // Default colors.
  private readonly colorGrey = colors.NEUTRAL_70;
  // Mouse-over colors.
  private readonly colorHandles = new three.Color('#eeeeee');
  private readonly colorHandlesHighlight = colors.WHITE;
  // Mouse-click colors.
  private readonly colorHandlesOnClick = new three.Color('#bbbbbb');
  private readonly orbitControls: IntrOrbitControls;
  // Raycaster to determine if mouse clicks or location falls on the widget's
  // control elements.
  private readonly raycaster = new three.Raycaster();
  private readonly interactiveObjects: three.Object3D[] = [];
  // Orthographic camera to render the widget.
  private readonly camera = new three.OrthographicCamera(-2, 2, 2, -2, 0, 4);
  // Size of the widget in pixels.
  private size: number = VIEW_WIDGET_BASE_PIXEL_SIZE;
  // Current quaternion slerp interpolation value (0.0 to 1.0)
  private currentInterpolationValue = 0.0;
  private time = performance.now();
  private readonly scene = new three.Scene();
  private readonly root = new three.Object3D();
  // Widget elements (axis and control points)
  private readonly posXAxisHelper: three.Sprite;
  private readonly posYAxisHelper: three.Sprite;
  private readonly posZAxisHelper: three.Sprite;
  private readonly negXAxisHelper: three.Sprite;
  private readonly negYAxisHelper: three.Sprite;
  private readonly negZAxisHelper: three.Sprite;
  // Used to suppress handling of mouse inside a handle when we are completing a click-drag.
  private draggingToRotate = false;

  private readonly quaternionUpdatedSubject = new Subject<three.Quaternion>();
  /** Emits the camera orientation requested by user interaction. */
  readonly quaternionUpdated: Observable<three.Quaternion> =
    this.quaternionUpdatedSubject;

  // Sprite currently highlighted by mouse-over.
  private highlitedSprite: three.Sprite | null = null;
  // Current highlight color.
  private highlightColor: three.Color = colors.X_AXIS;
  // Caches color updates to prevent redundant color changes.
  private needsColorUpdate = true;
  // Animation interpolation data.
  private animating = false;
  private readonly startQuaternion = new three.Quaternion();
  private readonly targetQuaternion = new three.Quaternion();

  private repaint = true;
  private readonly renderer: RendererClient;
  private readonly externalCameraSubscription: Subscription;
  private animationFrameHandle?: number;

  /**
   * @param canvas The canvas to render the widget into.
   * @param externalCamera The main camera whose orientation the widget mirrors.
   * @param worldRendererService Provides the renderer for the canvas.
   */
  constructor(
    private readonly canvas: HTMLCanvasElement,
    private readonly externalCamera: IntrinsicCamera,
    worldRendererService: WorldRendererService,
  ) {
    this.scene.add(this.root);
    this.camera.position.set(0, 0, 2);
    this.camera.up = new three.Vector3(0, 0, 1);
    // Create the widget's sprites and materials.
    const geometry = new three.BoxGeometry(0.8, 0.05, 0.05).translate(
      0.4,
      0,
      0,
    );

    const xAxis = new three.Mesh(
      geometry,
      new three.MeshBasicMaterial({
        color: colors.X_AXIS,
        toneMapped: false,
        depthWrite: false,
      }),
    );
    const yAxis = new three.Mesh(
      geometry,
      new three.MeshBasicMaterial({
        color: colors.Y_AXIS,
        toneMapped: false,
        depthWrite: false,
      }),
    );
    const zAxis = new three.Mesh(
      geometry,
      new three.MeshBasicMaterial({
        color: colors.Z_AXIS,
        toneMapped: false,
        depthWrite: false,
      }),
    );

    const negAxisMat = new three.MeshBasicMaterial({
      color: this.colorGrey,
      toneMapped: false,
    });
    const negXAxis = new three.Mesh(geometry, negAxisMat);
    const negYAxis = new three.Mesh(geometry, negAxisMat);
    const negZAxis = new three.Mesh(geometry, negAxisMat);

    negXAxis.rotation.z = Math.PI;
    negYAxis.rotation.z = Math.PI + Math.PI / 2;
    negZAxis.rotation.y = Math.PI + -Math.PI / 2;

    yAxis.rotation.z = Math.PI / 2;
    zAxis.rotation.y = -Math.PI / 2;

    this.root.add(xAxis);
    this.root.add(zAxis);
    this.root.add(yAxis);
    this.root.add(negXAxis);
    this.root.add(negYAxis);
    this.root.add(negZAxis);
    this.root.scale.setScalar(1.3);
    this.posXAxisHelper = new three.Sprite(
      createCircleSpriteMaterial(colors.X_AXIS, 'X'),
    );
    this.posXAxisHelper.scale.setScalar(VIEW_WIDGET_POS_AXIS_SCALE);
    this.posYAxisHelper = new three.Sprite(
      createCircleSpriteMaterial(colors.Y_AXIS, 'Y'),
    );
    this.posYAxisHelper.scale.setScalar(VIEW_WIDGET_POS_AXIS_SCALE);
    this.posZAxisHelper = new three.Sprite(
      createCircleSpriteMaterial(colors.Z_AXIS, 'Z'),
    );
    this.posZAxisHelper.scale.setScalar(VIEW_WIDGET_POS_AXIS_SCALE);
    this.negXAxisHelper = new three.Sprite(
      createCircleSpriteMaterial(this.colorGrey),
    );
    this.negXAxisHelper.scale.setScalar(VIEW_WIDGET_NEG_AXIS_SCALE);
    this.negYAxisHelper = new three.Sprite(
      createCircleSpriteMaterial(this.colorGrey),
    );
    this.negYAxisHelper.scale.setScalar(VIEW_WIDGET_NEG_AXIS_SCALE);
    this.negZAxisHelper = new three.Sprite(
      createCircleSpriteMaterial(this.colorGrey),
    );
    this.negZAxisHelper.scale.setScalar(VIEW_WIDGET_NEG_AXIS_SCALE);

    this.posXAxisHelper.position.x = 1;
    this.posYAxisHelper.position.y = 1;
    this.posZAxisHelper.position.z = 1;
    this.negXAxisHelper.position.x = -1;
    this.negYAxisHelper.position.y = -1;
    this.negZAxisHelper.position.z = -1;

    this.interactiveObjects.push(this.posXAxisHelper);
    this.interactiveObjects.push(this.posYAxisHelper);
    this.interactiveObjects.push(this.posZAxisHelper);
    this.interactiveObjects.push(this.negXAxisHelper);
    this.interactiveObjects.push(this.negYAxisHelper);
    this.interactiveObjects.push(this.negZAxisHelper);

    // Add all the sprites to the widget's render list.
    this.interactiveObjects.forEach((element) => {
      this.root.add(element);
    });

    this.externalCameraSubscription = this.externalCamera.updated.subscribe(
      () => {
        this.camera.quaternion.copy(this.externalCamera.quaternion);
        this.camera.position
          .copy(this.camera.getWorldDirection(new three.Vector3()))
          .normalize()
          .multiplyScalar(2)
          .negate();
        this.requestRepaint();
      },
    );

    canvas.addEventListener('contextmenu', this.handleContextMenu);
    canvas.addEventListener('pointerup', this.handleMouseClick);
    canvas.addEventListener('pointermove', this.handleMouseMove);
    this.orbitControls = new IntrOrbitControls(this.camera, canvas);
    this.orbitControls.addEventListener('change', () => {
      this.draggingToRotate = true;
      this.requestRepaint();
      this.quaternionUpdatedSubject.next(this.camera.quaternion.clone());
    });
    this.orbitControls.enableZoom = false;
    this.orbitControls.enablePan = false;
    this.orbitControls.target.copy(new three.Vector3(0, 0, 0));
    this.renderer = worldRendererService.createClient(canvas);
    this.renderer.setSize(
      VIEW_WIDGET_BASE_PIXEL_SIZE,
      VIEW_WIDGET_BASE_PIXEL_SIZE,
    );
    this.renderer.setClearColor(new three.Color(0, 0, 0), 0);
    this.mainLoop();
  }

  /** Stops rendering and releases the canvas event listeners. */
  dispose() {
    if (this.animationFrameHandle !== undefined) {
      cancelAnimationFrame(this.animationFrameHandle);
      this.animationFrameHandle = undefined;
    }
    this.externalCameraSubscription.unsubscribe();
    this.quaternionUpdatedSubject.complete();
    this.orbitControls.dispose();
    this.canvas.removeEventListener('contextmenu', this.handleContextMenu);
    this.canvas.removeEventListener('pointerup', this.handleMouseClick);
    this.canvas.removeEventListener('pointermove', this.handleMouseMove);
  }

  private readonly handleContextMenu = (e: Event) => {
    e.preventDefault();
  };

  private mainLoop() {
    this.update();
    this.render();
    this.animationFrameHandle = requestAnimationFrame(() => {
      this.mainLoop();
    });
  }

  private render() {
    if (!this.repaint) return;

    this.renderer.render(this.scene, this.camera);
    this.repaint = false;
  }

  /**
   * Render the widget's threejs elements.
   * White-circle background is rendered first.
   */
  requestRepaint() {
    this.repaint = true;
  }

  /**
   * Convert global mouse coordinates widget-space, NDC coordinates [-1 to 1]
   * for the widget's ortho camera. The widget-local coordinates are used to
   * detect mouse-overs and mouse-clicks hits. This function will not clip the
   * results. Values outside of -1 to 1 range, indicate the input x,y
   * coordinates do not fall within the widget's extents.
   * @param mousePosX Mouse x coordinate from js event.
   * @param mousePosY Mouse y coordinate from js event.
   * @returns Vector2 containing converted ndc coordinates of mouse cursor.
   */
  private mouseToWidgetCoords(
    mousePosX: number,
    mousePosY: number,
  ): three.Vector2 {
    // Flip Y coordinate for NDC
    return new three.Vector2(
      (mousePosX / this.size) * 2 - 1,
      -(mousePosY / this.size) * 2 + 1,
    );
  }

  /**
   * Mouse move handler used to detect mouseover events to highlight the element
   * under the mouse cursor. Actual element color changes are handled in the
   * update method for the widget.
   * This function also posts `mouse_in_view_widget` when the mouse is over
   * the widget, and `mouse_outside_view_widget` when the mouse is outside the
   * widget. These events are handled by viewController and used to allow
   * click-drag inside the view widget to operate the orbitCamera controls.
   * @param e Mouse event
   * @returns true if event is handled, false otherwise.
   */
  private readonly handleMouseMove = (e: MouseEvent): boolean => {
    const mousePosX = e.offsetX;
    const mousePosY = e.offsetY;
    const object = this.getSpriteObjectUnderMouse(mousePosX, mousePosY);
    let mouseoverSprite: three.Sprite | null = null;

    if (object != null) {
      this.highlightColor = this.colorHandlesHighlight;
      switch (object) {
        case this.posXAxisHelper:
          mouseoverSprite = this.posXAxisHelper;
          break;
        case this.posYAxisHelper:
          mouseoverSprite = this.posYAxisHelper;
          break;
        case this.posZAxisHelper:
          mouseoverSprite = this.posZAxisHelper;
          break;
        case this.negXAxisHelper:
          mouseoverSprite = this.negXAxisHelper;
          break;
        case this.negYAxisHelper:
          mouseoverSprite = this.negYAxisHelper;
          break;
        case this.negZAxisHelper:
          mouseoverSprite = this.negZAxisHelper;
          break;
        default:
          console.log(
            'Intersected unknown object type %d\n',
            object.userData['type'],
          );
          break;
      }
    }
    // Cache sprite to prevent redundant color changes on update.
    if (mouseoverSprite !== this.highlitedSprite) {
      this.highlitedSprite = mouseoverSprite;
      this.needsColorUpdate = true;
    }
    return true;
  };

  /**
   * Mouse click event handler. Used to trigger camera animation to reorient
   * the camera. Animation is handled by the widget's update method. This method
   * receives the target used by the frontend's orbit controls to maintain
   * the camera's target as the camera is animated.
   * @param e Mouse event.
   * @returns Returns true if the event was handled, false otherwise.
   */
  private readonly handleMouseClick = (e: MouseEvent): boolean => {
    if (e.button !== 0) {
      // Only handle left mouse button.
      return false;
    }
    const mousePosX = e.offsetX;
    const mousePosY = e.offsetY;
    // handleMouseClick is called on a mouseUp event

    let result = false;
    // Ignore mouse clicks we if are already animating, or
    // if we are moving the view.
    if (this.draggingToRotate) {
      this.draggingToRotate = false;
      return result;
    }
    if (this.animating === true) {
      return result;
    }
    const object = this.getSpriteObjectUnderMouse(mousePosX, mousePosY);
    if (object != null) {
      this.prepareAnimationData(object);
      this.animating = true;
      result = true;
      switch (object) {
        case this.posXAxisHelper:
          this.posXAxisHelper.material.color.copy(this.colorHandlesOnClick);
          break;
        case this.posYAxisHelper:
          this.posYAxisHelper.material.color.copy(this.colorHandlesOnClick);
          break;
        case this.posZAxisHelper:
          this.posZAxisHelper.material.color.copy(this.colorHandlesOnClick);
          break;
        case this.negXAxisHelper:
          this.negXAxisHelper.material.color.copy(this.colorHandlesOnClick);
          break;
        case this.negYAxisHelper:
          this.negYAxisHelper.material.color.copy(this.colorHandlesOnClick);
          break;
        case this.negZAxisHelper:
          this.negZAxisHelper.material.color.copy(this.colorHandlesOnClick);
          break;
        default:
          break;
      }
    }
    return result;
  };

  /**
   * Uses the threejs raycast to find the widget element
   * currently under the mouse cursor. Return null if no hits are
   * detected.
   * @param mousePosX Mouse x coordinate from js event.
   * @param mousePosY Mouse y coordinate from js event.
   * @returns Returns three.Object3D under the mouse cursor,
   *     otherwise null.
   */
  private getSpriteObjectUnderMouse(
    mousePosX: number,
    mousePosY: number,
  ): three.Object3D | null {
    // Convert mouse click coordinates into normalized coordinates
    // for ray intersections in the widget camera's orthographic space.
    const mouseInWidget = this.mouseToWidgetCoords(mousePosX, mousePosY);
    if (
      mouseInWidget.x < -1.0 ||
      mouseInWidget.x > 1.0 ||
      mouseInWidget.y < -1.0 ||
      mouseInWidget.y > 1.0
    ) {
      return null;
    }

    // Set up a raycaster.
    this.raycaster.setFromCamera(mouseInWidget, this.camera);

    // Cast rays into the widget's scene to see if they intersect any
    // of the widget's sprites.
    const intersections = this.raycaster.intersectObjects(
      this.interactiveObjects,
    );
    if (intersections.length === 0) {
      return null;
    }
    // We only care about the first (frontmost) intersection.
    return intersections[0].object;
  }

  /**
   * Sigmoid function to ease in and out of the quaternion animation.
   */
  sigmoid(x: number): number {
    return 1 / (1 + (x / (1 - x)) ** ANIMATION_SMOOTHING_VALUE);
  }
  /**
   * Update method for widget camera animations and color update.
   */
  update() {
    const now = performance.now();
    const delta = now - this.time;
    this.time = now;
    // If we are not animating, handle mouse-over color changes.
    this.size = VIEW_WIDGET_BASE_PIXEL_SIZE;

    if (!this.animating) {
      if (this.needsColorUpdate === true) {
        this.posXAxisHelper.material.color.set(this.colorHandles.getHex());
        this.posYAxisHelper.material.color.set(this.colorHandles.getHex());
        this.posZAxisHelper.material.color.set(this.colorHandles.getHex());
        this.negXAxisHelper.material.color.set(this.colorHandles.getHex());
        this.negYAxisHelper.material.color.set(this.colorHandles.getHex());
        this.negZAxisHelper.material.color.set(this.colorHandles.getHex());
        if (this.highlitedSprite) {
          this.highlitedSprite.material.color.set(this.highlightColor.getHex());
        }
        this.needsColorUpdate = false;
        this.repaint = true;
      }
      return;
    }
    this.repaint = true;
    // Animation is active. Handle quaternion interpolations.
    // Handle camera rotation animation updates.
    const step = delta * ANIMATION_INTERPOLATION_RATE;

    // Animate position by doing a slerp and then scaling the position on the
    // unit sphere..
    this.currentInterpolationValue += step;
    // If we reach 1.0, we are done animating.
    // Clamp interpolation value so we don't overshoot the target quaternion.
    if (this.currentInterpolationValue >= 1.0) {
      this.currentInterpolationValue = 1.0;
      this.animating = false;
    }

    const smoothInterpolation = this.sigmoid(this.currentInterpolationValue);
    // Apply quaternion slerp to position.
    const quaternion = new three.Quaternion();
    quaternion.slerpQuaternions(
      this.startQuaternion,
      this.targetQuaternion,
      smoothInterpolation,
    );

    this.quaternionUpdatedSubject.next(quaternion.normalize());
  }

  /**
   * Initializes animation data when the user clicks one of the axis controls.
   * @param object Object that was clicked by the user.
   */
  private prepareAnimationData(object: three.Object3D) {
    switch (object) {
      case this.posXAxisHelper:
        this.targetQuaternion.set(0.5, 0.5, 0.5, 0.5);
        this.targetQuaternion.normalize();
        break;

      case this.posYAxisHelper:
        this.targetQuaternion.set(0, 0.7071, 0.7071, 0);
        this.targetQuaternion.normalize();
        break;
      case this.posZAxisHelper:
        this.targetQuaternion.set(0, 0, 0, 1);
        this.targetQuaternion.normalize();
        break;
      case this.negXAxisHelper:
        this.targetQuaternion.set(-0.5, 0.5, 0.5, -0.5);
        this.targetQuaternion.normalize();
        break;
      case this.negYAxisHelper:
        this.targetQuaternion.set(0.7071, 0, 0, 0.7071);
        this.targetQuaternion.normalize();
        break;

      case this.negZAxisHelper:
        this.targetQuaternion.set(1, 0, 0, 0);
        this.targetQuaternion.normalize();
        break;
      default:
        console.error('ViewHelper: Invalid axis.');
        break;
    }

    this.currentInterpolationValue = 0.0;

    this.startQuaternion.copy(this.camera.quaternion);
  }
}
